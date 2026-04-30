package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/itoken417/go-simple-webapp/internal/member"
	"github.com/itoken417/go-simple-webapp/internal/secret"
)

func main() {
	key, err := secret.LoadKey()
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	members, err := member.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	nextID := nextID(members)
	scanner := bufio.NewScanner(os.Stdin)

	email := prompt(scanner, "メールアドレス: ")
	if email == "" {
		fmt.Fprintln(os.Stderr, "エラー: メールアドレスが空です")
		os.Exit(1)
	}
	for _, m := range members {
		if m.Email == email {
			fmt.Fprintf(os.Stderr, "エラー: %q は既に登録済みです\n", email)
			os.Exit(1)
		}
	}

	passwd := prompt(scanner, "パスワード: ")
	if passwd == "" {
		fmt.Fprintln(os.Stderr, "エラー: パスワードが空です")
		os.Exit(1)
	}
	if prompt(scanner, "パスワード（確認）: ") != passwd {
		fmt.Fprintln(os.Stderr, "エラー: パスワードが一致しません")
		os.Exit(1)
	}

	hashedB64, saltB64, err := member.HashPassword(passwd, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	id := strconv.Itoa(nextID)
	members = append(members, member.Member{
		ID:     id,
		Email:  email,
		Passwd: hashedB64,
		Salt:   saltB64,
	})

	if err := member.Save(members); err != nil {
		fmt.Fprintln(os.Stderr, "エラー:", err)
		os.Exit(1)
	}

	fmt.Printf("メンバーを追加しました: id=%s email=%s\n", id, email)
}

func nextID(members []member.Member) int {
	max := 0
	for _, m := range members {
		if n, err := strconv.Atoi(m.ID); err == nil && n > max {
			max = n
		}
	}
	return max + 1
}

func prompt(scanner *bufio.Scanner, label string) string {
	fmt.Print(label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
