package main

import (
	"crypto/md5"
	"encoding/binary"
)

func generateKey(longUrl string) string {
	return intToKey(checksum(longUrl)%(alphabetLength*alphabetLength)) + intToKey(getNextDatabaseId())
}

func checksum(str string) uint {
	sum := md5.Sum([]byte(str))
	return uint(binary.LittleEndian.Uint16(sum[:2])) // First 16 bits == first 2 bytes
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
