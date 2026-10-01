// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package fragment

import (
	"github.com/sagernet/sing-box/transport/finalmask"
	net "github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
)

func (c *Config) WrapConnClient(conn net.Conn, dest *net.Destination, dialer *finalmask.Dialer) (net.Conn, error) {
	return NewConnClient(c, conn, false)
}

func (c *Config) WrapConnServer(conn net.Conn) (net.Conn, error) {
	return NewConnServer(c, conn, true)
}
