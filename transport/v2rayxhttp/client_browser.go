// This file is derived from Xray-core (transport/internet/splithttp/browser_client.go),
// licensed under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Modifications for sing-box are licensed under GPL-3.0-or-later.

package v2rayxhttp

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/sagernet/sing-box/common/browserdialer"
	E "github.com/sagernet/sing/common/exceptions"
	M "github.com/sagernet/sing/common/metadata"
)

// BrowserDialerClient implements DialerClient in terms of the browser dialer.
// Like in Xray, only GET streams and packet uploads are possible, so the
// "auto" mode resolves to packet-up.
type BrowserDialerClient struct {
	transportConfig *Config
	dialer          *browserdialer.Dialer
}

func (c *BrowserDialerClient) IsClosed() bool {
	return false
}

func (c *BrowserDialerClient) OpenStream(ctx context.Context, url string, sessionId string, body io.Reader, uploadOnly bool) (io.ReadCloser, net.Addr, net.Addr, error) {
	if body != nil {
		return nil, nil, nil, E.New("bidirectional streaming for browser dialer not implemented yet")
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	c.transportConfig.FillStreamRequest(request, sessionId, "")

	conn, err := c.dialer.DialGet(ctx, request.URL.String(), request.Header, request.Cookies())
	if err != nil {
		return nil, M.Socksaddr{}, M.Socksaddr{}, err
	}
	return conn, conn.RemoteAddr(), conn.LocalAddr(), nil
}

func (c *BrowserDialerClient) PostPacket(ctx context.Context, url string, sessionId string, seqStr string, payload []byte) error {
	method := c.transportConfig.GetNormalizedUplinkHTTPMethod()
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	c.transportConfig.FillPacketRequest(request, sessionId, seqStr, payload)

	var bytes []byte
	if request.Body != nil {
		bytes, err = io.ReadAll(request.Body)
		if err != nil {
			return err
		}
	}

	return c.dialer.DialPacket(ctx, method, request.URL.String(), request.Header, request.Cookies(), bytes)
}
