package vlessenc

import (
	"net"
	"reflect"
	"unsafe"

	N "github.com/sagernet/sing/common/network"

	_ "github.com/sagernet/sing-vmess/vless"
)

//go:linkname visionTLSRegistry github.com/sagernet/sing-vmess/vless.tlsRegistry
var visionTLSRegistry []func(conn net.Conn) (loaded bool, netConn net.Conn, reflectType reflect.Type, reflectPointer uintptr)

// Make XTLS Vision work on top of VLESS Encryption, as in Xray
// (proxy.UnwrapRawConn and the CommonConn branches of the VLESS handlers):
// Vision drains CommonConn's own input/rawInput buffers and then copies
// directly to the connection below the encryption layer, which is the XorConn
// for "random" mode or the transport connection otherwise. Any TLS below
// CommonConn is intentionally not penetrated.
//
// The entry must come before sing-vmess' TLS entries, because those would
// otherwise find an outer TLS connection through CommonConn.Upstream().
func init() {
	visionTLSRegistry = append([]func(conn net.Conn) (bool, net.Conn, reflect.Type, uintptr){
		func(conn net.Conn) (bool, net.Conn, reflect.Type, uintptr) {
			commonConn, loaded := N.CastReader[*CommonConn](conn)
			if !loaded {
				return false, nil, nil, 0
			}
			return true, commonConn.Conn, reflect.TypeFor[CommonConn](), uintptr(unsafe.Pointer(commonConn))
		},
	}, visionTLSRegistry...)
}
