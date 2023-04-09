package main

import (
	"os"

	"github.com/Jenrykster/gotoh/app"
)

func main() {
	err := app.Init(os.Stdout, os.Stdin)
	if err != nil {
		os.Exit(1)
	}
}
