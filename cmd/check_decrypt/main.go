package main

import (
	"fmt"
	"os"

	"github.com/itoken417/go-simple-webapp/internal/secret"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	key, err := secret.LoadKey()
	if err != nil {
		fmt.Println("key error:", err)
		os.Exit(1)
	}
	enc := os.Getenv("SMTP_PASSWORD_ENCRYPTED")
	salt := os.Getenv("SMTP_PASSWORD_SALT")
	plain, err := secret.Decrypt(enc, salt, key)
	if err != nil {
		fmt.Println("decrypt error:", err)
		os.Exit(1)
	}
	fmt.Printf("password: %q\n", plain)
}
