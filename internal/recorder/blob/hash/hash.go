package hash

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/zeebo/blake3"
)

type Hash struct{
	hasher *blake3.Hasher
}

func NewHash() *Hash{
	return &Hash{
		hasher: blake3.New(),
	}
}

func (h *Hash) HashFile(path string) (string, error) {
    file, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer file.Close()

    if _, err := io.Copy(h.hasher, file); err != nil {
        return "", err
    }

    hash := h.hasher.Sum(nil)
    return hex.EncodeToString(hash), nil
}

func (h *Hash) HashContent(data []byte) string {
    hash := blake3.Sum256(data)
    return hex.EncodeToString(hash[:])
}

