// This file is derived from Xray-core (common/net), licensed under the
// Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
//
// Package xnet is the subset of Xray's common/net used by FinalMask.

package xnet

import (
	"encoding/binary"
	"strconv"

	errors "github.com/sagernet/sing-box/transport/finalmask/internal/xerrors"
)

// Port represents a network port in TCP and UDP protocol.
type Port uint16

// PortFromBytes converts a byte array to a Port, assuming bytes are in big endian order.
// @unsafe Caller must ensure that the byte array has at least 2 elements.
func PortFromBytes(port []byte) Port {
	return Port(binary.BigEndian.Uint16(port))
}

// PortFromInt converts an integer to a Port.
// @error when the integer is not positive or larger then 65535
func PortFromInt(val uint32) (Port, error) {
	if val > 65535 {
		return Port(0), errors.New("invalid port range: ", val)
	}
	return Port(val), nil
}

// PortFromString converts a string to a Port.
// @error when the string is not an integer or the integral value is a not a valid Port.
func PortFromString(s string) (Port, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return Port(0), errors.New("invalid port range: ", s)
	}
	return PortFromInt(uint32(val))
}

// Value return the corresponding uint16 value of a Port.
func (p Port) Value() uint16 {
	return uint16(p)
}

// String returns the string presentation of a Port.
func (p Port) String() string {
	return strconv.Itoa(int(p))
}

// FromPort returns the beginning port of this PortRange.

// ToPort returns the end port of this PortRange.

// Contains returns true if the given port is within the range of a PortRange.

// SinglePortRange returns a PortRange contains a single port.

type MemoryPortRange struct {
	From Port
	To   Port
}

func (r MemoryPortRange) Contains(port Port) bool {
	return r.From <= port && port <= r.To
}

type MemoryPortList []MemoryPortRange

func (mpl MemoryPortList) Contains(port Port) bool {
	for _, pr := range mpl {
		if pr.Contains(port) {
			return true
		}
	}
	return false
}
