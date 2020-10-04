package server

import (
	"crypto/rand"
	"encoding/base64"
)

func generateSlug(n int) string {
	raw := make([]byte, n)
	nonce := make([]byte, base64.StdEncoding.EncodedLen(len(raw)))
	if _, err := rand.Read(raw); err != nil {
		panic(err)
	}
	base64.StdEncoding.Encode(nonce, raw)
	return string(nonce)
}
