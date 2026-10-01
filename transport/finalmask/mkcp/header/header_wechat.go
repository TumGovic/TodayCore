// This file is derived from Xray-core (transport/internet/finalmask), licensed
// under the Mozilla Public License 2.0.
// Source: https://github.com/XTLS/Xray-core (v26.9.30)
// See LICENSES/MPL-2.0.txt and THIRD_PARTY_NOTICES.md.

package header

import "encoding/binary"

type wechat struct {
	sn uint32
}

func (*wechat) Size() int {
	return 13
}

func (h *wechat) Serialize(b []byte) {
	h.sn++
	b[0] = 0xa1
	b[1] = 0x08
	binary.BigEndian.PutUint32(b[2:], h.sn)
	b[6] = 0x00
	b[7] = 0x10
	b[8] = 0x11
	b[9] = 0x18
	b[10] = 0x30
	b[11] = 0x22
	b[12] = 0x30
}
