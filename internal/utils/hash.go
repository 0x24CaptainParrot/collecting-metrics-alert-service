package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/0x24CaptainParrot/collecting-metrics-alert-service.git/internal/models"
)

func ComputeSHA256(data interface{}, key string) (string, error) {
	h := hmac.New(sha256.New, []byte(key))

	switch v := data.(type) {
	case models.Metrics, []models.Metrics:
		jsonBytes, err := json.Marshal(data)
		if err != nil {
			return "", err
		}
		h.Write(jsonBytes)
	case string:
		h.Write([]byte(v))
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
