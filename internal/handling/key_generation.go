package handling

import (
	"hash/crc32"
)

// This alphabet will be used to generate the paths of the shortened URLs.
// It consists of the decimal digits and of the uppercase and lowercase letters, plus some special characters.
// Characters that could cause ambiguity or generate offensive words were removed.
const alphabet = "23456789BCDFGHJKLMNPQRSTVWXYZbcdfghjkmnpqrstvwxyz-_~!$&=@"
const alphabetLength = uint(len(alphabet))

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
