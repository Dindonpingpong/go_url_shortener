package secretary

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"

	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/config"
	"github.com/Dindonpingpong/yandex_practicum_go_url_shortener_service/service/secretary"
)

type Secretary struct {
	gcm   cipher.AEAD
	nonce []byte
}

var _ secretary.Secretary = (*Secretary)(nil)

func NewSecretaryService(c *config.SecretConfig) (*Secretary, error) {
	key := sha256.Sum256([]byte(c.UserKey))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := key[len(key)-gcm.NonceSize():]

	return &Secretary{
		gcm: gcm,
		nonce: nonce,
	}, nil
}

func (s *Secretary) Encode(data string) (string) {
	encryptedBytes := s.gcm.Seal(nil, s.nonce, []byte(data), nil)

	return hex.EncodeToString(encryptedBytes)
}

func (s *Secretary) Decode(msg string) (string, error) {
	msgBytes, err := hex.DecodeString(msg)

	if err != nil {
		return "", err
	}

	decryptedBytes, err := s.gcm.Open(nil, s.nonce, msgBytes, nil)

	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}
