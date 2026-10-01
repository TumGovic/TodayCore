//go:build !with_quic

package v2rayxhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/sagernet/sing-box/common/tls"
	C "github.com/sagernet/sing-box/constant"
)

func (c *Client) newHTTP3Transport(keepAlivePeriod time.Duration) http.RoundTripper {
	return nil
}

func newHTTP3Server(ctx context.Context, server *Server, tlsConfig tls.ServerConfig) (http3Server, error) {
	return nil, C.ErrQUICNotIncluded
}
