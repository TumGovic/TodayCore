// This file is derived from Xray-core (transport/internet/splithttp/hub.go),
// licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

package v2rayxhttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/common/tls"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"
	aTLS "github.com/sagernet/sing/common/tls"
	sHttp "github.com/sagernet/sing/protocol/http"

	"golang.org/x/net/http2"
)

var _ adapter.V2RayServerTransport = (*Server)(nil)

// Server is the XHTTP (SplitHTTP) server transport. It accepts every mode
// that Xray's server accepts (packet-up, stream-up and stream-one, unless
// restricted by "mode") over HTTP/1.1, h2 and h2c, and over HTTP/3 when the
// TLS ALPN is exactly ["h3"].
type Server struct {
	ctx        context.Context
	logger     logger.ContextLogger
	tlsConfig  tls.ServerConfig
	handler    adapter.V2RayServerTransportHandler
	config     *Config
	host       string
	path       string
	httpServer *http.Server
	h3Server   http3Server

	sessionAccess sync.Mutex
	sessions      sync.Map
}

type httpSession struct {
	uploadQueue *uploadQueue
	// for as long as the GET request is not opened by the client, this will be
	// open ("undone"), and the session may be expired within a certain TTL.
	// after the client connects, this becomes "done" and the session lives as
	// long as the GET request.
	isFullyConnected *doneInstance
}

func NewServer(ctx context.Context, logger logger.ContextLogger, options option.V2RayXHTTPOptions, tlsConfig tls.ServerConfig, handler adapter.V2RayServerTransportHandler) (adapter.V2RayServerTransport, error) {
	// downloadSettings only affects clients; Xray's server ignores it too.
	options.DownloadSettings = nil
	config, err := NewConfig(options)
	if err != nil {
		return nil, err
	}
	server := &Server{
		ctx:       ctx,
		logger:    logger,
		tlsConfig: tlsConfig,
		handler:   handler,
		config:    config,
		host:      config.Host,
		path:      config.GetNormalizedPath(),
	}
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(true)
	// sing-box TLS connections are not *tls.Conn, so net/http sees h2 over
	// TLS as prior-knowledge h2c. This also serves plaintext h2c, as Xray does.
	protocols.SetUnencryptedHTTP2(true)
	server.httpServer = &http.Server{
		Handler:           server,
		ReadHeaderTimeout: 4 * time.Second,
		MaxHeaderBytes:    config.GetNormalizedServerMaxHeaderBytes(),
		Protocols:         protocols,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			return log.ContextWithNewID(ctx)
		},
	}
	if isHTTP3ServerConfig(tlsConfig) {
		server.h3Server, err = newHTTP3Server(ctx, server, tlsConfig)
		if err != nil {
			return nil, err
		}
	}
	return server, nil
}

func isHTTP3ServerConfig(tlsConfig tls.ServerConfig) bool {
	if tlsConfig == nil {
		return false
	}
	nextProtos := tlsConfig.NextProtos()
	return len(nextProtos) == 1 && nextProtos[0] == "h3"
}

func (s *Server) upsertSession(sessionId string) *httpSession {
	// fast path
	currentSessionAny, ok := s.sessions.Load(sessionId)
	if ok {
		return currentSessionAny.(*httpSession)
	}

	// slow path
	s.sessionAccess.Lock()
	defer s.sessionAccess.Unlock()

	currentSessionAny, ok = s.sessions.Load(sessionId)
	if ok {
		return currentSessionAny.(*httpSession)
	}

	session := &httpSession{
		uploadQueue:      newUploadQueue(s.config.GetNormalizedScMaxBufferedPosts()),
		isFullyConnected: newDone(),
	}

	s.sessions.Store(sessionId, session)

	shouldReap := newDone()
	go func() {
		time.Sleep(30 * time.Second)
		shouldReap.Close()
	}()
	go func() {
		select {
		case <-shouldReap.Wait():
			s.sessions.Delete(sessionId)
			session.uploadQueue.Close()
		case <-session.isFullyConnected.Wait():
		}
	}()

	return session
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if len(s.host) > 0 && !isValidHTTPHost(request.Host, s.host) {
		s.logger.DebugContext(request.Context(), "failed to validate host, request: ", request.Host, ", config: ", s.host)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	if !strings.HasPrefix(request.URL.Path, s.path) {
		s.logger.DebugContext(request.Context(), "failed to validate path, request: ", request.URL.Path, ", config: ", s.path)
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	s.config.WriteResponseHeader(writer, request.Method, request.Header)
	length := int(s.config.GetNormalizedXPaddingBytes().rand())
	paddingConfig := XPaddingConfig{Length: length}

	if s.config.XPaddingObfsMode {
		paddingConfig.Placement = XPaddingPlacement{
			Placement: s.config.XPaddingPlacement,
			Key:       s.config.XPaddingKey,
			Header:    s.config.XPaddingHeader,
		}
		paddingConfig.Method = s.config.XPaddingMethod
	} else {
		paddingConfig.Placement = XPaddingPlacement{
			Placement: PlacementHeader,
			Header:    "X-Padding",
		}
	}

	s.config.ApplyXPaddingToResponse(writer, paddingConfig)

	if request.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		return
	}

	validRange := s.config.GetNormalizedXPaddingBytes()
	paddingValue, paddingPlacement := s.config.ExtractXPaddingFromRequest(request, s.config.XPaddingObfsMode)

	if !s.config.IsPaddingValid(paddingValue, validRange.From, validRange.To, s.config.XPaddingMethod) {
		s.logger.DebugContext(request.Context(), "invalid padding (", paddingPlacement, ") length: ", len(paddingValue))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	obfsPaddingAccepted := s.config.XPaddingObfsMode && paddingValue != ""

	sessionId, seqStr := s.config.ExtractMetaFromRequest(request, s.path)

	if sessionId == "" && s.config.Mode != "" && s.config.Mode != ModeAuto && s.config.Mode != ModeStreamOne && s.config.Mode != ModeStreamUp {
		s.logger.DebugContext(request.Context(), "stream-one mode is not allowed")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	source := sHttp.SourceAddress(request)

	var currentSession *httpSession
	if sessionId != "" {
		currentSession = s.upsertSession(sessionId)
	}
	scMaxEachPostBytes := int(s.config.GetNormalizedScMaxEachPostBytes().To)
	isUplinkRequest := false

	switch request.Method {
	case "GET":
		isUplinkRequest = seqStr != ""
	default:
		isUplinkRequest = true
	}

	uplinkDataKey := s.config.UplinkDataKey

	if isUplinkRequest && sessionId != "" { // stream-up, packet-up
		if seqStr == "" {
			if s.config.Mode != "" && s.config.Mode != ModeAuto && s.config.Mode != ModeStreamUp {
				s.logger.DebugContext(request.Context(), "stream-up mode is not allowed")
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			httpSC := &httpServerConn{
				doneInstance:   newDone(),
				Reader:         request.Body,
				ResponseWriter: writer,
			}
			err := currentSession.uploadQueue.Push(uploadPacket{
				Reader: httpSC,
			})
			if err != nil {
				s.logger.DebugContext(request.Context(), E.Cause(err, "failed to upload (PushReader)"))
				writer.WriteHeader(http.StatusConflict)
			} else {
				writer.Header().Set("X-Accel-Buffering", "no")
				writer.Header().Set("Cache-Control", "no-store")
				writer.WriteHeader(http.StatusOK)
				scStreamUpServerSecs := s.config.GetNormalizedScStreamUpServerSecs()
				hasLegacyRefererCompatMarker := request.Header.Get("Referer") != ""
				if (hasLegacyRefererCompatMarker || obfsPaddingAccepted) && scStreamUpServerSecs.To > 0 {
					go func() {
						for {
							_, err := httpSC.Write(bytes.Repeat([]byte{'X'}, int(s.config.GetNormalizedXPaddingBytes().rand())))
							if err != nil {
								break
							}
							time.Sleep(time.Duration(scStreamUpServerSecs.rand()) * time.Second)
						}
					}()
				}
				select {
				case <-request.Context().Done():
				case <-httpSC.Wait():
				}
			}
			httpSC.Close()
			return
		}

		if s.config.Mode != "" && s.config.Mode != ModeAuto && s.config.Mode != ModePacketUp {
			s.logger.DebugContext(request.Context(), "packet-up mode is not allowed")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		dataPlacement := s.config.GetNormalizedUplinkDataPlacement()
		var headerPayload []byte
		var err error
		if dataPlacement == PlacementAuto || dataPlacement == PlacementHeader {
			var headerPayloadChunks []string
			for i := 0; true; i++ {
				chunk := request.Header.Get(fmt.Sprintf("%s-%d", uplinkDataKey, i))
				if chunk == "" {
					break
				}
				headerPayloadChunks = append(headerPayloadChunks, chunk)
			}
			headerPayloadEncoded := strings.Join(headerPayloadChunks, "")
			headerPayload, err = base64.RawURLEncoding.DecodeString(headerPayloadEncoded)
			if err != nil {
				s.logger.DebugContext(request.Context(), "invalid base64 in header's payload: ", err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var cookiePayload []byte
		if dataPlacement == PlacementAuto || dataPlacement == PlacementCookie {
			var cookiePayloadChunks []string
			for i := 0; true; i++ {
				cookieName := fmt.Sprintf("%s_%d", uplinkDataKey, i)
				if c, _ := request.Cookie(cookieName); c != nil {
					cookiePayloadChunks = append(cookiePayloadChunks, c.Value)
				} else {
					break
				}
			}
			cookiePayloadEncoded := strings.Join(cookiePayloadChunks, "")
			cookiePayload, err = base64.RawURLEncoding.DecodeString(cookiePayloadEncoded)
			if err != nil {
				s.logger.DebugContext(request.Context(), "invalid base64 in cookies' payload: ", err)
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var bodyPayload []byte
		if dataPlacement == PlacementAuto || dataPlacement == PlacementBody {
			var readErr error
			if request.ContentLength > int64(scMaxEachPostBytes) {
				s.logger.DebugContext(request.Context(), "too large upload. scMaxEachPostBytes is set to ", scMaxEachPostBytes, " but request size exceed it. Adjust scMaxEachPostBytes on the server to be at least as large as client.")
				writer.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}
			if request.ContentLength > 0 {
				bodyPayload = make([]byte, request.ContentLength)
				_, readErr = io.ReadFull(request.Body, bodyPayload)
			} else {
				bodyPayload, readErr = io.ReadAll(io.LimitReader(request.Body, int64(scMaxEachPostBytes)+1))
			}
			if readErr != nil {
				s.logger.DebugContext(request.Context(), E.Cause(readErr, "failed to read body payload"))
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		var payload []byte
		switch dataPlacement {
		case PlacementHeader:
			payload = headerPayload
		case PlacementCookie:
			payload = cookiePayload
		case PlacementBody:
			payload = bodyPayload
		case PlacementAuto:
			payload = slices.Concat(headerPayload, cookiePayload, bodyPayload)
		}

		if len(payload) > scMaxEachPostBytes {
			s.logger.DebugContext(request.Context(), "too large upload. scMaxEachPostBytes is set to ", scMaxEachPostBytes, " but request size exceed it. Adjust scMaxEachPostBytes on the server to be at least as large as client.")
			writer.WriteHeader(http.StatusRequestEntityTooLarge)
			return
		}

		seq, err := strconv.ParseUint(seqStr, 10, 64)
		if err != nil {
			s.logger.DebugContext(request.Context(), E.Cause(err, "failed to upload (ParseUint)"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = currentSession.uploadQueue.Push(uploadPacket{
			Payload: payload,
			Seq:     seq,
		})
		if err != nil {
			s.logger.DebugContext(request.Context(), E.Cause(err, "failed to upload (PushPayload)"))
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if len(bodyPayload) == 0 {
			// Methods without a body are usually cached by default.
			writer.Header().Set("Cache-Control", "no-store")
		}

		writer.WriteHeader(http.StatusOK)
	} else if request.Method == "GET" || sessionId == "" { // stream-down, stream-one
		if sessionId != "" {
			// after GET is done, the connection is finished. disable automatic
			// session reaping, and handle it in defer
			currentSession.isFullyConnected.Close()
			defer s.sessions.Delete(sessionId)
		}

		// magic header instructs nginx + apache to not buffer response body
		writer.Header().Set("X-Accel-Buffering", "no")
		// A web-compliant header telling all middleboxes to disable caching.
		// Should be able to prevent overloading the cache, or stop CDNs from
		// teeing the response stream into their cache, causing slowdowns.
		writer.Header().Set("Cache-Control", "no-store")

		if !s.config.NoSSEHeader {
			// magic header to make the HTTP middle box consider this as SSE to disable buffer
			writer.Header().Set("Content-Type", "text/event-stream")
		}

		writer.WriteHeader(http.StatusOK)
		writer.(http.Flusher).Flush()

		httpSC := &httpServerConn{
			doneInstance:   newDone(),
			Reader:         request.Body,
			ResponseWriter: writer,
		}
		var localAddr net.Addr
		if la, ok := request.Context().Value(http.LocalAddrContextKey).(net.Addr); ok && la != nil {
			localAddr = la
		}
		conn := &splitConn{
			writer:     httpSC,
			reader:     httpSC,
			remoteAddr: source,
			localAddr:  localAddr,
		}
		if sessionId != "" { // if not stream-one
			conn.reader = currentSession.uploadQueue
		}

		s.handler.NewConnectionEx(request.Context(), conn, source, M.Socksaddr{}, N.OnceClose(func(it error) {
			conn.Close()
		}))

		// "A ResponseWriter may not be used after [Handler.ServeHTTP] has returned."
		select {
		case <-request.Context().Done():
		case <-httpSC.Wait():
		}

		conn.Close()
	} else {
		s.logger.DebugContext(request.Context(), "unsupported method: ", request.Method)
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type httpServerConn struct {
	sync.Mutex
	*doneInstance
	io.Reader // no need to Close request.Body
	http.ResponseWriter
}

func (c *httpServerConn) Write(b []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	if c.Done() {
		return 0, io.ErrClosedPipe
	}
	n, err := c.ResponseWriter.Write(b)
	if err == nil {
		c.ResponseWriter.(http.Flusher).Flush()
	}
	return n, err
}

func (c *httpServerConn) Close() error {
	c.Lock()
	defer c.Unlock()
	return c.doneInstance.Close()
}

func (s *Server) Network() []string {
	if s.h3Server != nil {
		return []string{N.NetworkUDP}
	}
	return []string{N.NetworkTCP}
}

func (s *Server) Serve(listener net.Listener) error {
	if s.h3Server != nil {
		return os.ErrInvalid
	}
	if s.tlsConfig != nil {
		if len(s.tlsConfig.NextProtos()) == 0 {
			s.tlsConfig.SetNextProtos([]string{http2.NextProtoTLS, "http/1.1"})
		}
		listener = aTLS.NewListener(listener, s.tlsConfig)
	}
	return s.httpServer.Serve(listener)
}

func (s *Server) ServePacket(listener net.PacketConn) error {
	if s.h3Server == nil {
		return os.ErrInvalid
	}
	return s.h3Server.Serve(listener)
}

func (s *Server) Close() error {
	return common.Close(common.PtrOrNil(s.httpServer), s.h3Server)
}
