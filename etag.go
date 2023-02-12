package optimisticlocker

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type EtagGenerator func(v any) (string, error)

var DefaultEtagGenerator = GenerateEtag

func GenerateEtag(v any) (string, error) {
	j, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(j)
	return hex.EncodeToString(hash[:]), nil
}
