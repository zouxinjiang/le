package lib

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
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

func Sha256(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}

func Hmac256(data []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func Hmac256X(data string, key string) string {
	return Hex(Hmac256([]byte(data), []byte(key)))
}
