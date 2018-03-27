// +build amd64

package checksum

//go:noescape
func sum(s uint, b []byte) uint
