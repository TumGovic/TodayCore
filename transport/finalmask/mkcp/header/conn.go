// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package header

import (
	"net"

	errors "github.com/sagernet/sing-box/transport/finalmask/internal/xerrors"
)

type headerConn struct {
	net.PacketConn
	header Header
}

func NewConnClient(c *Config, raw net.PacketConn) (net.PacketConn, error) {
	var header Header
	switch HeaderID(c.ID) {
	case DNS:
		var err error
		header, err = NewHeaderDNS(c.Domain)
		if err != nil {
			return nil, err
		}
	case DTLS:
		header = &dtls{}
	case SRTP:
		header = &srtp{}
	case UTP:
		header = &utp{}
	case WECHAT:
		header = &wechat{}
	case WIREGUARD:
		header = &wireguard{}
	default:
		return nil, errors.New("invalid id ", c.ID)
	}
	return &headerConn{
		PacketConn: raw,
		header:     header,
	}, nil
}

func NewConnServer(c *Config, raw net.PacketConn) (net.PacketConn, error) {
	return NewConnClient(c, raw)
}

func (c *headerConn) Size() int {
	return c.header.Size()
}

func (c *headerConn) ReadFrom(p []byte) (n int, addr net.Addr, err error) {
	return len(p) - c.header.Size(), nil, nil
}

func (c *headerConn) WriteTo(p []byte, addr net.Addr) (n int, err error) {
	c.header.Serialize(p)
	return len(p), nil
}
