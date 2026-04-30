package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/itoken417/go-simple-webapp/internal/member"
	"github.com/itoken417/go-simple-webapp/internal/secret"
)

func init() {
	register("02 メンバー追加", addMember)
}

func addMember(stdin *bufio.Reader) error {
	key, err := secret.LoadKey()
	if err != nil {
		return fmt.Errorf("鍵の読み込みに失敗: %w", err)
	}

	members, err := member.Load()
	if err != nil {
		return fmt.Errorf("メンバー読み込みに失敗: %w", err)
	}

	nextID := calcNextID(members)

	email, err := promptLine(stdin, "メールアドレス: ")
	if err != nil {
		return err
	}
	if email == "" {
		return fmt.Errorf("メールアドレスが空です")
	}
	for _, m := range members {
		if m.Email == email {
			return fmt.Errorf("%q は既に登録済みです", email)
		}
	}

	passwd, err := promptLine(stdin, "パスワード: ")
	if err != nil {
		return err
	}
	if passwd == "" {
		return fmt.Errorf("パスワードが空です")
	}
	confirm, err := promptLine(stdin, "パスワード（確認）: ")
	if err != nil {
		return err
	}
	if confirm != passwd {
		return fmt.Errorf("パスワードが一致しません")
	}

	hashedB64, saltB64, err := member.HashPassword(passwd, key)
	if err != nil {
		return fmt.Errorf("パスワードのハッシュ化に失敗: %w", err)
	}

	id := strconv.Itoa(nextID)
	members = append(members, member.Member{
		ID:     id,
		Email:  email,
		Passwd: hashedB64,
		Salt:   saltB64,
	})

	if err := member.Save(members); err != nil {
		return fmt.Errorf("メンバーの保存に失敗: %w", err)
	}

	fmt.Printf("メンバーを追加しました: id=%s email=%s\n", id, email)
	return nil
}

func calcNextID(members []member.Member) int {
	max := 0
	for _, m := range members {
		if n, err := strconv.Atoi(m.ID); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

func promptLine(r *bufio.Reader, label string) (string, error) {
	fmt.Print(label)
	input, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimRight(input, "\r\n"), nil
}
