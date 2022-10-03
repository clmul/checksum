package checksum

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"os"
	"strings"
	"testing"
)

func readFile(t testing.TB, f string) [][]byte {
	data, err := os.ReadFile(f)
	if err != nil {
		t.Fatal(err)
	}
	for data[len(data)-1] == '\n' {
		data = data[:len(data)-1]
	}
	var packets [][]byte
	for _, line := range strings.Split(string(data), "\n") {
		p, err := hex.DecodeString(line)
		if err != nil {
			t.Fatal(err)
		}
		packets = append(packets, p)
	}
	return packets
}

func test(t *testing.T, f string) {
	for i, packet0 := range readFile(t, f) {
		packet := make([]byte, len(packet0)+i)
		copy(packet[i:], packet0)
		Calc(packet[i:])
		if !bytes.Equal(packet[i:], packet0) {
			t.Log("got  ", packet)
			t.Log("want ", packet0)
			t.Fatal("wrong checksum")
		}
	}
}

func TestCalcICMP(t *testing.T) {
	test(t, "testdata/icmp")
}

func TestCalcTCP(t *testing.T) {
	test(t, "testdata/tcp")
}

func TestCalcUDP(t *testing.T) {
	test(t, "testdata/udp")
}

func TestCalcIPv4(t *testing.T) {
	lines := readFile(t, "testdata/tcp")
	for _, packet0 := range lines {
		packet := make([]byte, len(packet0))
		copy(packet, packet0)
		CalcIPv4(packet)

		if !bytes.Equal(packet, packet0) {
			t.Log("got  ", packet)
			t.Log("want ", packet0)
			t.Fatal("wrong checksum")
		}
	}
}

func TestUpdateByte(t *testing.T) {
	lines := readFile(t, "testdata/tcp")
	for _, packet0 := range lines {
		packet := make([]byte, len(packet0))
		copy(packet, packet0)

		position := rand.Int()%8 + 1
		v := byte(rand.Int())

		packet0[position] = v
		Calc(packet0)

		UpdateByte(packet, position, v)
		if !bytes.Equal(packet, packet0) {
			t.Log("got  ", packet)
			t.Log("want ", packet0)
			t.Fatalf("wrong checksum, position: %v, v: %v", position, v)
		}
	}
}

func BenchmarkCalc(b *testing.B) {
	lines := readFile(b, "testdata/tcp")
	var packets [][]byte
	for _, p := range lines {
		if len(p) > 1024 {
			packets = append(packets, p[:1024])
		}
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pi := i % len(packets)
		p := packets[pi]
		Calc(p)
		b.SetBytes(int64(len(p)))
	}
}
func BenchmarkCalcIPv4(b *testing.B) {
	packets := readFile(b, "testdata/tcp")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pi := i % len(packets)
		p := packets[pi]
		CalcIPv4(p)
		b.SetBytes(int64(len(p)))
	}
}
func BenchmarkUpdateByte(b *testing.B) {
	packets := readFile(b, "testdata/tcp")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pi := i % len(packets)
		p := packets[pi]
		UpdateByte(p, 8, 254)
		b.SetBytes(int64(len(p)))
	}
}

func BenchmarkSum(b *testing.B) {
	const size = 1023
	p := make([]byte, size)
	for i := range p {
		p[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum(0, p)
		b.SetBytes(size)
	}
}
