//go:build amd64 || arm64

package checksum

//go:noescape
func sum(s uint64, b []byte) uint64
