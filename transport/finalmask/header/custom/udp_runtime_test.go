// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package custom

import (
	"testing"
)

func TestDSLUDPServerMatchCapturesSavedValues(t *testing.T) {
	conn, err := NewConnServerUDP(&UDPConfig{
		Client: []*UDPItem{
			{
				Rand: 2,
				Save: "txid",
			},
			{
				Var: "txid",
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	server := conn.(*udpCustomServerConn)
	if !server.header.Match([]byte{0x01, 0x02, 0x01, 0x02}, nil) {
		t.Fatal("expected packet to match")
	}

	if got := string(server.header.vars["txid"]); got != string([]byte{0x01, 0x02}) {
		t.Fatalf("unexpected saved txid: %x", server.header.vars["txid"])
	}
}

func TestDSLUDPServerRejectsMalformedVarReference(t *testing.T) {
	conn, err := NewConnServerUDP(&UDPConfig{
		Client: []*UDPItem{
			{
				Rand: 2,
				Save: "txid",
			},
			{
				Var: "txid",
			},
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	server := conn.(*udpCustomServerConn)
	if server.header.Match([]byte{0x01, 0x02, 0x03, 0x04}, nil) {
		t.Fatal("expected packet mismatch")
	}
}
