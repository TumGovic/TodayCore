// This file is derived from Xray-core (common/net), licensed under the
// Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Package xnet is the subset of Xray's common/net used by FinalMask.

package xnet

// PacketConnWrapper wraps a PacketConn into a Conn with a fixed destination address.
type PacketConnWrapper struct {
	PacketConn
	Dest Addr
}

func (c *PacketConnWrapper) Read(p []byte) (int, error) {
	n, _, err := c.PacketConn.ReadFrom(p)
	return n, err
}

func (c *PacketConnWrapper) Write(p []byte) (int, error) {
	return c.PacketConn.WriteTo(p, c.Dest)
}

func (c *PacketConnWrapper) RemoteAddr() Addr {
	return c.Dest
}
