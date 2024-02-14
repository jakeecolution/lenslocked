package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nRead, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, fmt.Errorf("error generating random bytes: %w", err)
	}
	if n < nRead {
		return nil, fmt.Errorf("did not read correct number of bytes: want %d, got %d", n, nRead)
	}
	return b, nil
}

func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func HashToken(token string) (string, error) {
	return "", nil
}
