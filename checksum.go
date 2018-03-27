package checksum

import (
	"encoding/binary"
)

const (
	ICMP = 1
	TCP  = 6
	UDP  = 17

	ProtocolOffset = 9

	IPv4ChecksumOffset = 10
	ICMPChecksumOffset = 2
	TCPChecksumOffset  = 16
	UDPChecksumOffset  = 6
)

func CalcIPv4(packet []byte) {
	ihl := uint16(packet[0] & 0xf)
	ipv4HeaderLen := ihl * 4
	update(packet[:ipv4HeaderLen], packet, IPv4ChecksumOffset)
}

func Calc(packet []byte) {
	ihl := uint16(packet[0] & 0xf)
	protocol := packet[ProtocolOffset]
	ipv4HeaderLen := ihl * 4

	var l4ChecksumI int
	var h []byte
	var moreFragments bool
	var fragmentOffset uint16

	switch protocol {
	case UDP:
		l4ChecksumI = UDPChecksumOffset
	case TCP:
		l4ChecksumI = TCPChecksumOffset
	case ICMP:
		l4ChecksumI = ICMPChecksumOffset
	default:
		goto bypassL4
	}

	moreFragments = packet[6]&0x20 != 0
	fragmentOffset = uint16(packet[6]&0x1f)<<8 | uint16(packet[7])
	if moreFragments || fragmentOffset != 0 {
		goto bypassL4
	}

	switch protocol {
	case TCP, UDP:
		h = make([]byte, 12)
		copy(h[0:8], packet[12:20])
		h[8] = 0
		h[9] = protocol
		binary.BigEndian.PutUint16(h[10:12], uint16(len(packet))-ipv4HeaderLen)
	}
	updateWithHeader(h, packet[ipv4HeaderLen:], packet, int(ipv4HeaderLen)+l4ChecksumI)

bypassL4:
	update(packet[:ipv4HeaderLen], packet, IPv4ChecksumOffset)
}

func updateWithHeader(h, b, packet []byte, i int) {
	binary.LittleEndian.PutUint16(packet[i:], 0)
	s := sum(sum(0, h), b)
	for s>>16 > 0 {
		s = s&0xffff + s>>16
	}
	binary.LittleEndian.PutUint16(packet[i:], uint16(^s))
}

func update(b, packet []byte, i int) {
	binary.LittleEndian.PutUint16(packet[i:], 0)
	s := sum(0, b)
	for s>>16 > 0 {
		s = s&0xffff + s>>16
	}
	binary.LittleEndian.PutUint16(packet[i:], uint16(^s))
}
