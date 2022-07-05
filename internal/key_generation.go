package main

import (
	"hash/crc32"
)

func generateKey(longUrl string, id uint) string {
	prefix := intToKey(checksum(longUrl) % (alphabetLength * alphabetLength))
	if len(prefix) == 1 {
		prefix = alphabet[0:1] + prefix
	}
	return prefix + intToKey(id)
}

func checksum(str string) uint {
	return uint(crc32.ChecksumIEEE([]byte(str)))
}

func intToKey(num uint) string {
	if num == 0 {
		return alphabet[0:1]
	}
	var intToKeyRec func(uint) []byte
	intToKeyRec = func(num uint) []byte {
		if num == 0 {
			return []byte{}
		}
		return append(intToKeyRec(num/alphabetLength), alphabet[num%alphabetLength])
	}
	return string(intToKeyRec(num))
}
