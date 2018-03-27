package checksum

import (
	"encoding/binary"
)

// update byte in IP header
func UpdateByte(packet []byte, position int, v byte) {
	old := packet[position]
	if old == v {
		return
	}

	var before, after uint16
	if position&1 == 0 {
		before = binary.LittleEndian.Uint16(packet[position:])
		packet[position] = v
		after = binary.LittleEndian.Uint16(packet[position:])
	} else {
		before = binary.LittleEndian.Uint16(packet[position-1:])
		packet[position] = v
		after = binary.LittleEndian.Uint16(packet[position-1:])
	}

	sum := uint32(^binary.LittleEndian.Uint16(packet[IPv4ChecksumOffset:]))
	sum += uint32(^before)
	sum += uint32(after)
	sum = sum&0xffff + sum>>16
	binary.LittleEndian.PutUint16(packet[IPv4ChecksumOffset:], uint16(^sum))
}
