package cryptography

import (
	"fmt"

	"crypto/aes"
	"crypto/cipher"
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
	salt, err = c.generateRandomBytes(c.rtConfig.GetArgon2SaltLength())
	if err != nil {
		return key, salt, err
	}
	key = argon2.IDKey([]byte(password), salt,
		c.rtConfig.GetArgon2Time(),
		c.rtConfig.GetArgon2Memory(),
		c.rtConfig.GetArgon2Threads(),
		c.rtConfig.GetArgon2KeyLength(),
	)

	return key, salt, nil
}

func (c *Cryptographer) VerifyKeyAndGetDEK(password string, salt, dekEnc []byte) (dek []byte, err error) {
	key := argon2.IDKey([]byte(password), salt,
		c.rtConfig.GetArgon2Time(),
		c.rtConfig.GetArgon2Memory(),
		c.rtConfig.GetArgon2Threads(),
		c.rtConfig.GetArgon2KeyLength(),
	)

	aesgcm, err := c.createNewAESGCM(key)
	if err != nil {
		return nil, fmt.Errorf("can't create new aesgcm: %w", err)
	}

	dek, err = aesgcm.Open(nil, nil, dekEnc, nil)
	if err != nil {
		return nil, fmt.Errorf("can't decrypt dek: %w", err)
	}

	return dek, nil
}

func (c *Cryptographer) EncryptData(data string, dek []byte) ([]byte, error) {
	aesgcm, err := c.createNewAESGCM(dek)
	if err != nil {
		return nil, fmt.Errorf("can't create new aesgcm: %w", err)
	}

	encData := aesgcm.Seal(nil, nil, []byte(data), nil)

	return encData, nil
}

func (c *Cryptographer) DecryptData(data []byte, dek []byte) (string, error) {
	aesgcm, err := c.createNewAESGCM(dek)
	if err != nil {
		return "", fmt.Errorf("can't create new aesgcm: %w", err)
	}

	plaintext, err := aesgcm.Open(nil, nil, data, nil)
	if err != nil {
		return "", fmt.Errorf("can't decrypt data: %w", err)
	}

	return string(plaintext), nil
}

func (c *Cryptographer) generateRandomBytes(length uint32) ([]byte, error) {
	if length < 8 || length > 64 {
		return nil, fmt.Errorf("invalid length")
	}
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error while generating random bytes: %w", err)
	}

	return salt, nil
}

func (c *Cryptographer) createNewAESGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("can't create aes block: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("can't create aesgcm: %w", err)
	}

	return aesgcm, nil
}
