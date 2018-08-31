// +build !amd64,!arm64

package checksum

import (
	"encoding/binary"
)

func sum(s uint64, b []byte) uint64 {
	i := len(b)
	if i&1 != 0 {
		i -= 1
		v := uint64(b[i])
		s += v
		if s < v {
			s++
		}
	}
	for i&6 != 0 {
		i -= 2
		v := uint64(binary.LittleEndian.Uint16(b[i:]))
		s += v
		if s < v {
			s++
		}
	}
	for i > 0 {
		i -= 8
		v := binary.LittleEndian.Uint64(b[i:])
		s += v
		if s < v {
			s++
		}
	}
	return s
}
