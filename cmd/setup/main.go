package main

import (
	"bufio"
	"fmt"
	"os"
)

type step struct {
	name string
	fn   func(*bufio.Reader) error
}

var steps []step

func register(name string, fn func(*bufio.Reader) error) {
	steps = append(steps, step{name: name, fn: fn})
}

func main() {
	stdin := bufio.NewReader(os.Stdin)

	for _, s := range steps {
		fmt.Printf("\n=== %s ===\n", s.name)
		if err := s.fn(stdin); err != nil {
			fmt.Fprintln(os.Stderr, "エラー:", err)
			os.Exit(1)
		}
	}
}
