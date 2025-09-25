package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
)

type Hasher struct {
	h hash.Hash
}

func (h *Hasher) Sign(data []byte) (string, error) {
	h.h.Reset()

	_, err := h.h.Write(data)
	if err != nil {
		return "", err
	}
	hashBytes := h.h.Sum(nil)
	return hex.EncodeToString(hashBytes), nil
}

func (h *Hasher) VerifyHex(data []byte, expectedHash string) bool {
	actualHash, err := h.Sign(data)
	if err != nil {
		return false
	}
	return hmac.Equal([]byte(actualHash), []byte(expectedHash))
}

func NewHash(key []byte) *Hasher {
	h := hmac.New(sha256.New, key)

	return &Hasher{
		h: h,
	}
}
