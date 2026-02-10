package myjwt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
)

type Manager struct {
	aead cipher.AEAD
}

// key 必须是 16 / 24 / 32 字节
func NewManager(key string) (*Manager, error) {
	realKey := deriveKey(key) // 自动变成32字节
	block, err := aes.NewCipher(realKey)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &Manager{
		aead: aead,
	}, nil
}

func (m *Manager) CreateToken(v any) (string, error) {
	// 1. 序列化 JSON
	plain, err := json.Marshal(v)
	if err != nil {
		return "", err
	}

	// 2. 生成 nonce (12字节)
	nonce := make([]byte, m.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 3. 加密（自带认证tag）
	cipherText := m.aead.Seal(nil, nonce, plain, nil)

	// 4. 拼接 nonce + 密文
	final := append(nonce, cipherText...)

	// 5. base64url
	return base64.RawURLEncoding.EncodeToString(final), nil
}

func (m *Manager) ParseToken(token string, out any) error {

	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return errors.New("invalid base64")
	}

	nonceSize := m.aead.NonceSize()
	if len(raw) < nonceSize {
		return errors.New("invalid token")
	}

	nonce := raw[:nonceSize]
	cipherText := raw[nonceSize:]

	// 解密（自动校验篡改）
	plain, err := m.aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return errors.New("invalid token or tampered")
	}

	return json.Unmarshal(plain, out)
}

func deriveKey(keyStr string) []byte {
	hash := sha256.Sum256([]byte(keyStr))
	return hash[:]
}
