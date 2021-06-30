package crypto

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

func Sha1(origin string) string {
	h := sha1.New()
	h.Write([]byte(origin))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Sha256(origin string) string {
	h := sha256.New()
	h.Write([]byte(origin))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMacSha512(secret, origin string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(origin))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func HMacSha1(secret, origin string) string {
	h := hmac.New(sha512.New, []byte(secret))
	h.Write([]byte(origin))
	return fmt.Sprintf("%x", h.Sum(nil))
}
