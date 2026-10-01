// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package custom

import (
	"github.com/sagernet/sing-box/transport/finalmask"
	net "github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
)

func (c *TCPConfig) WrapConnClient(conn net.Conn, dest *net.Destination, dialer *finalmask.Dialer) (net.Conn, error) {
	return NewConnClientTCP(c, conn)
}

func (c *TCPConfig) WrapConnServer(conn net.Conn) (net.Conn, error) {
	return NewConnServerTCP(c, conn)
}

func (c *UDPConfig) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewConnClientUDP(c, conn)
}

func (c *UDPConfig) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return NewConnServerUDP(c, conn)
}

func (c *UDPStandaloneConfig) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewConnClientUDPStandalone(c, conn)
}

func (c *UDPStandaloneConfig) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return NewConnServerUDPStandalone(c, conn)
}
