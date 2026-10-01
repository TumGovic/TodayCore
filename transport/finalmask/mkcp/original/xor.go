// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

//go:build !amd64
// +build !amd64

package original

// xorfwd performs XOR forwards in words, x[i] ^= x[i-4], i from 0 to len
func xorfwd(x []byte) {
	for i := 4; i < len(x); i++ {
		x[i] ^= x[i-4]
	}
}

// xorbkd performs XOR backwords in words, x[i] ^= x[i-4], i from len to 0
func xorbkd(x []byte) {
	for i := len(x) - 1; i >= 4; i-- {
		x[i] ^= x[i-4]
	}
}
