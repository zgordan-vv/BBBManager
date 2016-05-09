package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

func CW(n int) string {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
