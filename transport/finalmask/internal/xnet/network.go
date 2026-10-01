// This file is derived from Xray-core (common/net), licensed under the
// Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Package xnet is the subset of Xray's common/net used by FinalMask.

package xnet

import "net"

func (n Network) SystemString() string {
	switch n {
	case Network_TCP:
		return "tcp"
	case Network_UDP:
		return "udp"
	case Network_UNIX:
		return "unix"
	default:
		return "unknown"
	}
}

// HasNetwork returns true if the network list has a certain network.
func HasNetwork(list []Network, network Network) bool {
	for _, value := range list {
		if value == network {
			return true
		}
	}
	return false
}

// Network mirrors Xray's protobuf Network enum.
type Network int32

const (
	Network_Unknown Network = 0
	Network_TCP     Network = 2
	Network_UDP     Network = 3
	Network_UNIX    Network = 4
)

// UDPAddrFromAddrPort is net.UDPAddrFromAddrPort.
var UDPAddrFromAddrPort = net.UDPAddrFromAddrPort

var (
	ListenPacket = net.ListenPacket
	Pipe         = net.Pipe
	IPv4         = net.IPv4
	ErrClosed    = net.ErrClosed
)
