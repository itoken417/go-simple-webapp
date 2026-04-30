package secret

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/itoken417/goutils/logger"
)

const SecretFile = ".secret"

// LoadKey は APP_SECRET_KEY を環境変数または .secret ファイルから読み込む。
// 環境変数が設定されている場合はそちらを優先する（本番環境向け）。
func LoadKey() ([]byte, error) {
	if v := os.Getenv("APP_SECRET_KEY"); v != "" {
		return parseHexKey(v)
	}
	f, err := os.Open(SecretFile)
	if err != nil {
		err := fmt.Errorf("APP_SECRET_KEY が未設定で %s も見つかりません", SecretFile)
		logger.Log(err)
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if after, ok := strings.CutPrefix(scanner.Text(), "APP_SECRET_KEY="); ok {
			return parseHexKey(after)
		}
	}
	err = fmt.Errorf("%s に APP_SECRET_KEY が見つかりません", SecretFile)
	logger.Log(err)
	return nil, err
}

// GenerateKeyFile は新しい 32 バイト鍵を生成して .secret ファイルに書き込む。
func GenerateKeyFile() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		err := fmt.Errorf("鍵生成失敗: %w", err)
		logger.Log(err)
		return nil, err
	}
	f, err := os.OpenFile(SecretFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		err := fmt.Errorf("%s 作成失敗: %w", SecretFile, err)
		logger.Log(err)
		return nil, err
	}
	defer f.Close()
	if _, err := fmt.Fprintf(f, "APP_SECRET_KEY=%s\n", hex.EncodeToString(key)); err != nil {
		err := fmt.Errorf("鍵書き込み失敗: %w", err)
		logger.Log(err)
		return nil, err
	}
	return key, nil
}

func parseHexKey(s string) ([]byte, error) {
	key, err := hex.DecodeString(strings.TrimSpace(s))
	if err != nil {
		err := fmt.Errorf("鍵のデコード失敗: %w", err)
		logger.Log(err)
		return nil, err
	}
	if len(key) != 32 {
		err := fmt.Errorf("鍵は 32 バイト必要です（現在 %d バイト）", len(key))
		logger.Log(err)
		return nil, err
	}
	return key, nil
}
