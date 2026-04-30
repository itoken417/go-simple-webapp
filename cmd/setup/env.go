package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/itoken417/go-simple-webapp/internal/secret"
	"github.com/itoken417/goutils/logger"
)

const (
	exampleFile = ".env.example"
	envFile     = ".env"
)

func main() {
	logger.Init(true, false)
	defer logger.Destory()

	stdin := bufio.NewReader(os.Stdin)

	if _, err := os.Stat(envFile); err == nil {
		fmt.Printf("%s はすでに存在します。上書きしますか？ [y/N]: ", envFile)
		answer, _ := stdin.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(answer)) != "y" {
			fmt.Println("キャンセルしました。")
			return
		}
	}

	encKey, err := loadOrGenerateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "鍵の準備に失敗しました: %v\n", err)
		os.Exit(1)
	}

	src, err := os.Open(exampleFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s を開けません: %v\n", exampleFile, err)
		os.Exit(1)
	}
	defer src.Close()

	fmt.Println("各項目の値を入力してください（Enter でデフォルト値を使用）。")
	fmt.Println()

	var lines []string
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			lines = append(lines, line)
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			lines = append(lines, line)
			continue
		}
		name, defaultVal := parts[0], parts[1]

		fmt.Printf("%s [%s]: ", name, defaultVal)
		input, _ := stdin.ReadString('\n')
		val := strings.TrimRight(input, "\r\n")
		if val == "" {
			val = defaultVal
		}

		// SMTP_PASSWORD は暗号化して 2 行に展開する
		if name == "SMTP_PASSWORD" {
			ciphertext, nonce, err := secret.Encrypt(val, encKey)
			if err != nil {
				fmt.Fprintf(os.Stderr, "暗号化失敗: %v\n", err)
				os.Exit(1)
			}
			lines = append(lines, "SMTP_PASSWORD_ENCRYPTED="+ciphertext)
			lines = append(lines, "SMTP_PASSWORD_SALT="+nonce)
			continue
		}

		lines = append(lines, name+"="+val)
	}

	dst, err := os.Create(envFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s を作成できません: %v\n", envFile, err)
		os.Exit(1)
	}
	defer dst.Close()

	w := bufio.NewWriter(dst)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "書き込みエラー: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n%s を作成しました。\n", envFile)
}

func loadOrGenerateKey() ([]byte, error) {
	if _, err := os.Stat(secret.SecretFile); os.IsNotExist(err) {
		key, err := secret.GenerateKeyFile()
		if err != nil {
			return nil, err
		}
		fmt.Printf("%s を新規作成しました（権限 600）\n", secret.SecretFile)
		return key, nil
	}
	key, err := secret.LoadKey()
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s から既存の鍵を読み込みました\n", secret.SecretFile)
	return key, nil
}
