package main

import (
	"os"

	"github.com/Jenrykster/gotoh/api"
	"github.com/Jenrykster/gotoh/app"
)

func main() {
	anilistClient := api.NewAnilistClient()
	err := app.Init(os.Stdout, os.Stdin, &anilistClient)
	if err != nil {
		os.Exit(1)
	}
}
