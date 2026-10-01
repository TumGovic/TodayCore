// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package header

type Header interface {
	Size() int
	Serialize(b []byte)
}

type HeaderID int

const (
	DNS HeaderID = iota
	DTLS
	SRTP
	UTP
	WECHAT
	WIREGUARD
)
