package secret

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/itoken417/goutils/logger"
)

// Encrypt は AES-256-GCM でプレーンテキストを暗号化する。
// 返り値は base64 エンコードされた ciphertext と nonce（salt）。
func Encrypt(plaintext string, key []byte) (ciphertextB64, nonceB64 string, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		err := fmt.Errorf("cipher 作成失敗: %w", err)
		logger.Log(err)
		return "", "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err := fmt.Errorf("GCM 作成失敗: %w", err)
		logger.Log(err)
		return "", "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		err := fmt.Errorf("nonce 生成失敗: %w", err)
		logger.Log(err)
		return "", "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext),
		base64.StdEncoding.EncodeToString(nonce),
		nil
}

// Decrypt は AES-256-GCM で復号する。
func Decrypt(ciphertextB64, nonceB64 string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		err := fmt.Errorf("ciphertext デコード失敗: %w", err)
		logger.Log(err)
		return "", err
	}
	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		err := fmt.Errorf("nonce デコード失敗: %w", err)
		logger.Log(err)
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		err := fmt.Errorf("cipher 作成失敗: %w", err)
		logger.Log(err)
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		err := fmt.Errorf("GCM 作成失敗: %w", err)
		logger.Log(err)
		return "", err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		err := fmt.Errorf("復号失敗: %w", err)
		logger.Log(err)
		return "", err
	}
	return string(plaintext), nil
}
