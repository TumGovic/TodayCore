//go:build with_quic

// This file is derived from Xray-core (transport/internet/splithttp/dialer.go
// and hub.go), licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.
//
// XHTTP over HTTP/3. Xray uses the apernet quic-go fork, sing-box uses the
// sagernet one; both speak standard QUIC v1 + HTTP/3 on the wire, and both
// parrot Chrome's QUIC Initial when ChromeParrot is set.

package v2rayxhttp

import (
	"context"
	stdTLS "crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/sagernet/quic-go"
	"github.com/sagernet/quic-go/http3"
	"github.com/sagernet/sing-box/common/tls"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/log"
	"github.com/sagernet/sing-quic"
	congestion_meta2 "github.com/sagernet/sing-quic/congestion_meta2"
	E "github.com/sagernet/sing/common/exceptions"
	N "github.com/sagernet/sing/common/network"
)

const quicgoH3KeepAlivePeriod = 10 * time.Second

func useBBR(conn *quic.Conn) {
	conn.SetCongestionControl(congestion_meta2.NewBbrSenderWithProfile(conn.InitialPacketSize(), congestion_meta2.ProfileStandard))
}

func (c *Client) newHTTP3Transport(keepAlivePeriod time.Duration) http.RoundTripper {
	quicConfig := &quic.Config{
		MaxIdleTimeout:          connIdleTimeout,
		KeepAlivePeriod:         quicgoH3KeepAlivePeriod,
		MaxIncomingStreams:      -1,
		DisablePathMTUDiscovery: !C.IsLinux && !C.IsWindows && !C.IsDarwin,
		ChromeParrot:            true,
	}
	if keepAlivePeriod > 0 {
		quicConfig.KeepAlivePeriod = keepAlivePeriod
	} else if keepAlivePeriod < 0 {
		quicConfig.KeepAlivePeriod = 0
	}
	return &http3.Transport{
		QUICConfig:      quicConfig,
		TLSClientConfig: &stdTLS.Config{},
		Dial: func(ctx context.Context, addr string, _ *stdTLS.Config, quicConfig *quic.Config) (*quic.Conn, error) {
			udpConn, err := c.dialer.DialContext(ctx, N.NetworkUDP, c.serverAddr)
			if err != nil {
				return nil, err
			}
			quicConn, err := qtls.DialEarly(ctx, udpConn, c.tlsConfig, quicConfig)
			if err != nil {
				udpConn.Close()
				return nil, err
			}
			context.AfterFunc(quicConn.Context(), func() { udpConn.Close() })
			useBBR(quicConn)
			return quicConn, nil
		},
	}
}

type http3ServerImpl struct {
	ctx       context.Context
	server    *Server
	tlsConfig tls.ServerConfig
	h3Server  *http3.Server
	listener  qtls.EarlyListener
}

func newHTTP3Server(ctx context.Context, server *Server, tlsConfig tls.ServerConfig) (http3Server, error) {
	err := qtls.ConfigureHTTP3(tlsConfig)
	if err != nil {
		return nil, err
	}
	return &http3ServerImpl{
		ctx:       ctx,
		server:    server,
		tlsConfig: tlsConfig,
		h3Server: &http3.Server{
			Handler: server,
			ConnContext: func(ctx context.Context, conn *quic.Conn) context.Context {
				useBBR(conn)
				return log.ContextWithNewID(ctx)
			},
		},
	}, nil
}

func (s *http3ServerImpl) Serve(packetConn net.PacketConn) error {
	listener, err := qtls.ListenEarly(packetConn, s.tlsConfig, &quic.Config{
		MaxIdleTimeout:          connIdleTimeout,
		DisablePathMTUDiscovery: !C.IsLinux && !C.IsWindows && !C.IsDarwin,
	})
	if err != nil {
		return err
	}
	s.listener = listener
	err = s.h3Server.ServeListener(listener)
	if err != nil && !E.IsClosedOrCanceled(err) {
		return err
	}
	return nil
}

func (s *http3ServerImpl) Close() error {
	err := s.h3Server.Close()
	if s.listener != nil {
		return E.Errors(err, s.listener.Close())
	}
	return err
}
