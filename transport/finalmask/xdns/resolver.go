// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package xdns

import (
	"errors"
	"net"

	"github.com/sagernet/sing-box/transport/finalmask"
)

type Resolver interface {
	Addr() *net.UDPAddr
	Read(p []byte) (int, error)
	Send(p []byte)
	Close()
}

// NewResolver builds a resolver from a *TCPResolverProto or *UDPResolverProto
// (a serial.TypedMessage in Xray).
func NewResolver(config any, dialer *finalmask.Dialer) (Resolver, error) {
	switch v := config.(type) {
	case *TCPResolverProto:
		return NewTCPResolver(v, dialer)
	case *UDPResolverProto:
		return NewUDPResolver(v, dialer)
	default:
		return nil, errors.New("unknown proto")
	}
}
