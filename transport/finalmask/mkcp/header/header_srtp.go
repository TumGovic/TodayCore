// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package header

import "encoding/binary"

type srtp struct {
	header uint16
	number uint16
}

func (*srtp) Size() int {
	return 4
}

func (h *srtp) Serialize(b []byte) {
	h.number++
	binary.BigEndian.PutUint16(b, h.header)
	binary.BigEndian.PutUint16(b[2:], h.number)
}
