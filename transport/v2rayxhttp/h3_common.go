package v2rayxhttp

import "net"

type http3Server interface {
	Serve(packetConn net.PacketConn) error
	Close() error
}
