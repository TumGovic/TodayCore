// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package salamander

import (
	"github.com/sagernet/sing-box/transport/finalmask"
	net "github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
)

func (c *Config) HeaderConn() {}

func (c *Config) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewSalamanderConnClient(c, conn)
}

func (c *Config) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return NewSalamanderConnServer(c, conn)
}

func (c *GeckoConfig) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewGeckoConnClient(c, conn)
}

func (c *GeckoConfig) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return NewGeckoConnServer(c, conn)
}
