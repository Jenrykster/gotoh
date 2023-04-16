package main

import (
	"os"

	"github.com/Jenrykster/gotoh/api"
	"github.com/Jenrykster/gotoh/app"
	"github.com/Jenrykster/gotoh/auth"
	"github.com/Jenrykster/gotoh/utils"
)

var env = utils.GetEnv()

func main() {
	if len(env.USER_TOKEN) == 0 {
		auth.GetAnilistUserOAuthToken()
	} else {
		anilistClient := api.NewAnilistClient(env.USER_TOKEN)
		err := app.Init(os.Stdout, os.Stdin, &anilistClient)
		if err != nil {
			os.Exit(1)
		}
	}
}
