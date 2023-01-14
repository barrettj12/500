package util

import "fmt"

// Utility functions

// E evaluates the given function and panics if it returns a non-nil error.
func E0(err error) {
	if err != nil {
		panic(err)
	}
}

func E[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// Coloured printing

func Red(s string) string {
	return fmt.Sprintf("\u001b[31m%s\u001b[0m", s)
}

func Grey(s string) string {
	return fmt.Sprintf("\033[38;2;200;200;200m%s\033[0m", s)
}
