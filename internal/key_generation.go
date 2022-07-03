package main

import (
	"hash/crc32"
)

func generateKey(longUrl string) string {
	return intToKey(checksum(longUrl)%(alphabetLength*alphabetLength)) + intToKey(getNextDatabaseId())
}

func checksum(str string) uint {
	return uint(crc32.ChecksumIEEE([]byte(str)))
}

func intToKey(num uint) string {
	var intToKeyRec func(uint) []byte
	intToKeyRec = func(num uint) []byte {
		if num == 0 {
			return []byte{}
		}
		return append(intToKeyRec(num/alphabetLength), alphabet[num%alphabetLength])
	}
	return string(intToKeyRec(num))
}
