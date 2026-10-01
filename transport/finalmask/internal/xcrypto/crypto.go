// This file is derived from Xray-core (common/crypto), licensed under the
// Mozilla Public License 2.0.

// Package xcrypto is the subset of Xray's common/crypto used by FinalMask.
package xcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"math/big"

	common "github.com/sagernet/sing-box/transport/finalmask/internal/xcommon"
)

// [,)
func RandBetween(from int64, to int64) int64 {
	if from > to {
		from, to = to, from
	}
	if d := to - from; d == 0 || d == 1 {
		return from
	}
	bigInt, _ := rand.Int(rand.Reader, big.NewInt(to-from))
	return from + bigInt.Int64()
}

// [,]
func RandBytesBetween(b []byte, from, to byte) {
	common.Must2(rand.Read(b))

	if from > to {
		from, to = to, from
	}

	if to-from == 255 {
		return
	}

	for i := range b {
		b[i] = from + b[i]%(to-from+1)
	}
}

func NewAesGcm(key []byte) cipher.AEAD {
	block := common.Must2(aes.NewCipher(key))
	aead := common.Must2(cipher.NewGCM(block))
	return aead
}
