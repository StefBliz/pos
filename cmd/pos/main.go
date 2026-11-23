package main

import (
	"fmt"
	"os"

	"pos/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка приложения:", err)
		os.Exit(1)
	}
}
