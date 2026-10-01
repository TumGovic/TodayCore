package v2rayxhttp

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/sagernet/sing-box/common/tls"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/logger"
	M "github.com/sagernet/sing/common/metadata"
	N "github.com/sagernet/sing/common/network"

	"github.com/stretchr/testify/require"
)

type echoHandler struct{}

func (echoHandler) NewConnectionEx(ctx context.Context, conn net.Conn, source M.Socksaddr, destination M.Socksaddr, onClose N.CloseHandlerFunc) {
	go func() {
		_, err := io.Copy(conn, conn)
		conn.Close()
		if onClose != nil {
			onClose(err)
		}
	}()
}

func newTestCertificate(t *testing.T) (string, string) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "example.com"},
		DNSNames:     []string{"example.com"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	require.NoError(t, err)
	keyDER, err := x509.MarshalECPrivateKey(key)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})),
		string(pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}))
}

func testXHTTP(t *testing.T, alpn []string, options option.V2RayXHTTPOptions) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var serverTLS tls.ServerConfig
	var clientTLS tls.Config
	if alpn != nil {
		certificate, key := newTestCertificate(t)
		var err error
		serverTLS, err = tls.NewServer(ctx, logger.NOP(), option.InboundTLSOptions{
			Enabled:     true,
			ALPN:        alpn,
			Certificate: []string{certificate},
			Key:         []string{key},
		})
		require.NoError(t, err)
		require.NoError(t, serverTLS.Start())
		defer serverTLS.Close()
		clientTLS, err = tls.NewClient(ctx, logger.NOP(), "example.com", option.OutboundTLSOptions{
			Enabled:    true,
			ServerName: "example.com",
			Insecure:   true,
			ALPN:       alpn,
		})
		require.NoError(t, err)
	}

	server, err := NewServer(ctx, logger.NOP(), options, serverTLS, echoHandler{})
	require.NoError(t, err)
	defer server.Close()

	var serverAddr M.Socksaddr
	if len(alpn) == 1 && alpn[0] == "h3" {
		require.Equal(t, []string{N.NetworkUDP}, server.Network())
		packetConn, err := net.ListenPacket("udp", "127.0.0.1:0")
		require.NoError(t, err)
		defer packetConn.Close()
		go server.ServePacket(packetConn)
		serverAddr = M.SocksaddrFromNet(packetConn.LocalAddr())
	} else {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		defer listener.Close()
		go server.Serve(listener)
		serverAddr = M.SocksaddrFromNet(listener.Addr())
	}

	client, err := NewClient(ctx, N.SystemDialer, serverAddr, options, clientTLS)
	require.NoError(t, err)
	defer client.Close()

	for i := 0; i < 2; i++ {
		conn, err := client.DialContext(ctx)
		require.NoError(t, err)
		payload := make([]byte, 256*1024)
		rand.Read(payload)
		go func() {
			conn.Write(payload[:1000])
			conn.Write(payload[1000:])
		}()
		received := make([]byte, len(payload))
		done := make(chan error, 1)
		go func() {
			_, err := io.ReadFull(conn, received)
			done <- err
		}()
		select {
		case err = <-done:
			require.NoError(t, err)
		case <-time.After(15 * time.Second):
			t.Fatal("timeout")
		}
		require.True(t, bytes.Equal(payload, received))
		conn.Close()
	}
}

func TestXHTTPServer(t *testing.T) {
	transports := map[string][]string{
		"http1.1": nil,
		"tls-h1":  {"http/1.1"},
		"tls-h2":  {"h2"},
	}
	if C.WithQUIC {
		transports["h3"] = []string{"h3"}
	}
	for name, alpn := range transports {
		for _, mode := range []string{ModePacketUp, ModeStreamUp, ModeStreamOne} {
			t.Run(name+"/"+mode, func(t *testing.T) {
				testXHTTP(t, alpn, option.V2RayXHTTPOptions{
					Host: "example.com",
					Path: "/xhttp",
					Mode: mode,
				})
			})
		}
	}
}

func TestXHTTPServerObfs(t *testing.T) {
	testXHTTP(t, []string{"h2"}, option.V2RayXHTTPOptions{
		Path:                "/xhttp?ed=2560",
		Mode:                ModePacketUp,
		XPaddingObfsMode:    true,
		XPaddingPlacement:   PlacementCookie,
		XPaddingMethod:      PaddingMethodTokenish,
		XPaddingKey:         "token",
		SessionIDPlacement:  PlacementHeader,
		SeqPlacement:        PlacementQuery,
		UplinkDataPlacement: PlacementHeader,
		UplinkHTTPMethod:    "PUT",
		// Header-carried uplink needs small posts and a larger header limit,
		// exactly as with Xray.
		ScMaxEachPostBytes:   option.XHTTPRange{From: 4000, To: 4000},
		ServerMaxHeaderBytes: 16384,
	})
}
