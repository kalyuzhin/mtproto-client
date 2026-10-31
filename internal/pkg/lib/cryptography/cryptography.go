package cryptography

import (
	"fmt"

	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

type rtConfig interface {
	GetArgon2Memory() uint32
	GetArgon2Threads() uint8
	GetArgon2Time() uint32
	GetArgon2SaltLength() uint32
	GetArgon2KeyLength() uint32
}

type Cryptographer struct {
	rtConfig rtConfig
}

// GenerateArgon2IDKey – generates argon2id key
func (c *Cryptographer) GenerateArgon2IDKey(password string) (key, salt []byte, err error) {
	salt, err = c.generateSalt(c.rtConfig.GetArgon2SaltLength())
	if err != nil {
		return key, salt, err
	}
	key = argon2.IDKey([]byte(password), salt,
		c.rtConfig.GetArgon2Time(),
		c.rtConfig.GetArgon2Memory(),
		c.rtConfig.GetArgon2Threads(),
		c.rtConfig.GetArgon2KeyLength())

	return key, salt, nil
}

func (c *Cryptographer) generateSalt(saltLength uint32) ([]byte, error) {
	if saltLength < 8 || saltLength > 32 {
		return nil, fmt.Errorf("invalid saltLength")
	}
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error while generating salt: %w", err)
	}

	return salt, nil
}
