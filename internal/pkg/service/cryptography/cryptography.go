package cryptography

import (
	"fmt"

	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type rtConfig interface {
}

type Service struct {
	rtConfig rtConfig
}

// GenerateArgon2IDKey – generates argon2id key
func (s *Service) GenerateArgon2IDKey(password string) ([]byte, error) {
	salt, err := s.generateSalt(16)
	if err != nil {
		return nil, err
	}
	argon2.IDKey([]byte(password), salt, 15, 1, 2, 32)

	return nil, nil
}

func (s *Service) generateSalt(saltLength int8) ([]byte, error) {
	salt := make([]byte, 0, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error while generating salt: %w", err)
	}

	return salt, nil
}
