//go:build !amd64 && !arm64

package checksum

import (
	"encoding/binary"
	"math/bits"
)

func sum(s uint64, b []byte) uint64 {
	i := len(b)
	var carry uint64
	if i&1 != 0 {
		i -= 1
		v := uint64(b[i])
		s, carry = bits.Add64(s, v, 0)
		s += carry
	}
	for i&6 != 0 {
		i -= 2
		v := uint64(binary.LittleEndian.Uint16(b[i:]))
		s, carry = bits.Add64(s, v, 0)
		s += carry
	}
	for i > 0 {
		i -= 8
		v := binary.LittleEndian.Uint64(b[i:])
		s, carry = bits.Add64(s, v, 0)
		s += carry
	}
	return s
}
