// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package sudoku

import (
	"github.com/sagernet/sing-box/transport/finalmask"
	net "github.com/sagernet/sing-box/transport/finalmask/internal/xnet"
)

// Sudoku in finalmask mode is a pure appearance transform with no standalone handshake.
// TCP always keeps classic sudoku on uplink and uses packed downlink optimization on server writes.
func (c *Config) WrapConnClient(conn net.Conn, dest *net.Destination, dialer *finalmask.Dialer) (net.Conn, error) {
	return newPackedDirectionalConn(conn, c, true)
}

func (c *Config) WrapConnServer(conn net.Conn) (net.Conn, error) {
	return newPackedDirectionalConn(conn, c, false)
}

func newPackedDirectionalConn(raw net.Conn, config *Config, readPacked bool) (net.Conn, error) {
	pureReader, pureWriter, err := newPureReaderWriter(raw, config)
	if err != nil {
		return nil, err
	}
	packedReader, packedWriter, err := newPackedReaderWriter(raw, config)
	if err != nil {
		return nil, err
	}

	reader, writer := pureReader, pureWriter
	if readPacked {
		reader = packedReader
	} else {
		writer = packedWriter
	}

	return newWrappedConn(raw, reader, writer), nil
}

func (c *Config) WrapPacketConnClient(conn net.PacketConn, dest *net.Destination, dialer *finalmask.Dialer) (net.PacketConn, error) {
	return NewUDPConn(conn, c)
}

func (c *Config) WrapPacketConnServer(conn net.PacketConn, addr net.Addr, lc *finalmask.ListenConfig) (net.PacketConn, error) {
	return NewUDPConn(conn, c)
}
