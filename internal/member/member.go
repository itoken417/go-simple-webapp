package member

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

const (
	dataFile     = "data/member.json"
	pbkdf2Iter   = 600_000
	pbkdf2KeyLen = 32
	saltLen      = 16
)

type Member struct {
	ID     string `json:"id"`
	Email  string `json:"email"`
	Passwd string `json:"passwd"`
	Salt   string `json:"salt"`
}

// HashPassword はパスワードをPBKDF2-SHA256でハッシュ化する。
// key はアプリ共通のペッパーとして使用する（APP_SECRET_KEY）。
func HashPassword(passwd string, key []byte) (hashedB64, saltB64 string, err error) {
	salt := make([]byte, saltLen)
	if _, err = rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("salt生成失敗: %w", err)
	}
	hashed, err := pbkdf2.Key(sha256.New, passwd+string(key), salt, pbkdf2Iter, pbkdf2KeyLen)
	if err != nil {
		return "", "", fmt.Errorf("ハッシュ化失敗: %w", err)
	}
	return base64.StdEncoding.EncodeToString(hashed),
		base64.StdEncoding.EncodeToString(salt),
		nil
}

// VerifyPassword はパスワードが正しいか検証する。
// key はハッシュ化時と同じペッパー（APP_SECRET_KEY）を渡す。
func VerifyPassword(passwd, hashedB64, saltB64 string, key []byte) (bool, error) {
	salt, err := base64.StdEncoding.DecodeString(saltB64)
	if err != nil {
		return false, fmt.Errorf("saltデコード失敗: %w", err)
	}
	expected, err := base64.StdEncoding.DecodeString(hashedB64)
	if err != nil {
		return false, fmt.Errorf("passwdデコード失敗: %w", err)
	}
	hashed, err := pbkdf2.Key(sha256.New, passwd+string(key), salt, pbkdf2Iter, pbkdf2KeyLen)
	if err != nil {
		return false, fmt.Errorf("ハッシュ化失敗: %w", err)
	}
	return subtle.ConstantTimeCompare(hashed, expected) == 1, nil
}

// Load はdata/member.jsonからメンバー一覧を読み込む。
func Load() ([]Member, error) {
	f, err := os.Open(dataFile)
	if os.IsNotExist(err) {
		return []Member{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("member.json読み込み失敗: %w", err)
	}
	defer f.Close()
	var members []Member
	if err := json.NewDecoder(f).Decode(&members); err != nil {
		return nil, fmt.Errorf("member.jsonパース失敗: %w", err)
	}
	return members, nil
}

// Save はメンバー一覧をdata/member.jsonに書き込む。
func Save(members []Member) error {
	if err := os.MkdirAll("data", 0700); err != nil {
		return fmt.Errorf("dataディレクトリ作成失敗: %w", err)
	}
	f, err := os.OpenFile(dataFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("member.json書き込み失敗: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(members)
}
