package lib

import (
	"crypto/sha1"
	"fmt"
)

func Sha1(data []byte) []byte {
	sha := sha1.New()
	sha.Write(data)
	return sha.Sum(nil)
}

func Sha1String(data string) string {
	return string(Sha1([]byte(data)))
}

func Hex(data []byte) string {
	return fmt.Sprintf("%x", data)
}

func Sha1HexString(data string) string {
	return Hex(Sha1([]byte(data)))
}
