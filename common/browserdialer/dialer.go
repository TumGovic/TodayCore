// This file is derived from Xray-core (transport/internet/browser_dialer),
// licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

// Package browserdialer is Xray's Browser Dialer: a local web page, opened in
// a real browser, performs the XHTTP / WebSocket requests with the browser's
// own TLS and HTTP stack, so the traffic carries a genuine browser
// fingerprint. The page talks to us over WebSocket connections to
// /websocket?token=<csrf token>; every connection carries one task.
//
// Unlike Xray, which enables a single global dialer through the
// XRAY_BROWSER_DIALER environment variable, the listen address is configured
// per transport. Like in Xray, a dialer lives for the whole process and is
// shared by every transport using the same address: sing-box closes
// transports to reset their connections and keeps using them afterwards.
package browserdialer

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	E "github.com/sagernet/sing/common/exceptions"
	"github.com/sagernet/sing/common/logger"

	"github.com/coder/websocket"
	"github.com/gofrs/uuid/v5"
)

//go:embed dialer.html
var webpage []byte

type task struct {
	Method         string `json:"method"`
	URL            string `json:"url"`
	Extra          any    `json:"extra,omitempty"`
	StreamResponse bool   `json:"streamResponse"`
}

type Dialer struct {
	conns chan *websocket.Conn
	ctx   context.Context
}

var (
	access  sync.Mutex
	dialers = make(map[string]*Dialer)
)

// Get returns the browser dialer listening on address, starting it if
// necessary.
func Get(address string, logger logger.Logger) (*Dialer, error) {
	access.Lock()
	defer access.Unlock()
	if dialer, loaded := dialers[address]; loaded {
		return dialer, nil
	}
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, E.Cause(err, "listen browser dialer")
	}
	csrfToken := uuid.Must(uuid.NewV4()).String()
	page := bytes.ReplaceAll(webpage, []byte("csrfToken"), []byte(csrfToken))
	dialer := &Dialer{
		conns: make(chan *websocket.Conn, 256),
		ctx:   context.Background(),
	}
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/websocket" {
				if r.URL.Query().Get("token") != csrfToken {
					return
				}
				conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{InsecureSkipVerify: true})
				if err != nil {
					logger.Error(E.Cause(err, "browser dialer: upgrade"))
					return
				}
				conn.SetReadLimit(-1)
				select {
				case dialer.conns <- conn:
				default:
					conn.CloseNow()
				}
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Write(page)
		}),
		ReadHeaderTimeout: 4 * time.Second,
	}
	go server.Serve(listener)
	logger.Info("browser dialer listening on ", listener.Addr(), ", open it in a browser to start dialing")
	dialers[address] = dialer
	return dialer, nil
}

// DialWS asks the browser to open a WebSocket to uri, sending earlyData as
// the Sec-WebSocket-Protocol header like Xray does. The returned connection
// carries the WebSocket payload as a byte stream.
func (d *Dialer) DialWS(ctx context.Context, uri string, earlyData []byte) (net.Conn, error) {
	conn, err := d.dialTask(ctx, task{
		Method:         "WS",
		URL:            uri,
		StreamResponse: true,
		Extra: struct {
			Protocol string `json:"protocol,omitempty"`
		}{base64.RawURLEncoding.EncodeToString(earlyData)},
	})
	if err != nil {
		return nil, err
	}
	return d.netConn(conn), nil
}

type httpExtra struct {
	Referrer string            `json:"referrer,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
	Cookies  map[string]string `json:"cookies,omitempty"`
}

func httpExtraFromHeadersAndCookies(headers http.Header, cookies []*http.Cookie) *httpExtra {
	if len(headers) == 0 {
		return nil
	}

	extra := httpExtra{}
	if referrer := headers.Get("Referer"); referrer != "" {
		extra.Referrer = referrer
		headers.Del("Referer")
	}

	if len(headers) > 0 {
		extra.Headers = make(map[string]string)
		for header := range headers {
			extra.Headers[header] = headers.Get(header)
		}
	}

	if len(cookies) > 0 {
		extra.Cookies = make(map[string]string)
		for _, cookie := range cookies {
			extra.Cookies[cookie.Name] = cookie.Value
		}
	}

	return &extra
}

// DialGet asks the browser to fetch uri and stream the response body back.
func (d *Dialer) DialGet(ctx context.Context, uri string, headers http.Header, cookies []*http.Cookie) (net.Conn, error) {
	conn, err := d.dialTask(ctx, task{
		Method:         "GET",
		URL:            uri,
		Extra:          httpExtraFromHeadersAndCookies(headers, cookies),
		StreamResponse: true,
	})
	if err != nil {
		return nil, err
	}
	return d.netConn(conn), nil
}

// DialPacket asks the browser to send one request with payload as the body.
func (d *Dialer) DialPacket(ctx context.Context, method string, uri string, headers http.Header, cookies []*http.Cookie, payload []byte) error {
	conn, err := d.dialTask(ctx, task{
		Method:         method,
		URL:            uri,
		Extra:          httpExtraFromHeadersAndCookies(headers, cookies),
		StreamResponse: false,
	})
	if err != nil {
		return err
	}
	defer conn.CloseNow()
	err = conn.Write(ctx, websocket.MessageBinary, payload)
	if err != nil {
		return err
	}
	return checkOK(ctx, conn)
}

func (d *Dialer) dialTask(ctx context.Context, task task) (*websocket.Conn, error) {
	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	for {
		var conn *websocket.Conn
		select {
		case conn = <-d.conns:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		if conn.Write(ctx, websocket.MessageText, data) != nil {
			conn.CloseNow()
			continue
		}
		err = checkOK(ctx, conn)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
}

func checkOK(ctx context.Context, conn *websocket.Conn) error {
	_, message, err := conn.Read(ctx)
	if err != nil {
		conn.CloseNow()
		return err
	}
	if s := string(message); s != "ok" {
		conn.CloseNow()
		return E.New(s)
	}
	return nil
}

func (d *Dialer) netConn(conn *websocket.Conn) net.Conn {
	return websocket.NetConn(d.ctx, conn, websocket.MessageBinary)
}
