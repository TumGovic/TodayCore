// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package udphop

import (
	"github.com/sagernet/sing-box/transport/finalmask"
	errors "github.com/sagernet/sing-box/transport/finalmask/internal/xerrors"
	net "github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
)

func (c *Config) HandleDial() {}

func (c *Config) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewUDPHopConn(c, dest, dialer)
}

func (c *Config) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return nil, errors.New("udphop: client only")
}
