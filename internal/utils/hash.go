package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5String returns the MD5 hash of the given string.
func MD5String(s string) string {
	return MD5Bytes([]byte(s))
}

// MD5Bytes returns the MD5 hash of the given bytes.
func MD5Bytes(b []byte) string {
	h := md5.New()
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}
