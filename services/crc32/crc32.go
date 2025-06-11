package crc32

import (
	"fmt"
	"hash/crc32"
)

func CreateChecksum(input string) string {
	crc32q := crc32.MakeTable(crc32.Castagnoli)
	checksum := crc32.Checksum([]byte(input), crc32q)
	return fmt.Sprintf("%08x", checksum)
}
