// Package internal for project
package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Verification struct
type Verification struct {
	secret []byte
}

// NewVerification return Verification
func NewVerification(secret string) *Verification {
	return &Verification{
		secret: []byte(secret),
	}
}

// CheckToken func
func (v *Verification) CheckToken(token string, p string) bool {
	return token == v.GenerateToken(p)
}

// GenerateToken func
func (v *Verification) GenerateToken(p string) string {
	mac := hmac.New(sha256.New, v.secret)
	_, _ = mac.Write([]byte(p))
	expectedMAC := mac.Sum(nil)

	return toHex(expectedMAC)
}

// toHex return hex string
func toHex(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)

	return string(dst)
}
