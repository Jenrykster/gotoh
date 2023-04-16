package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jenrykster/gotoh/utils"
)

const tokenRequest = "https://anilist.co/api/v2/oauth/authorize?client_id=%s&response_type=code"

var env = utils.GetEnv()

func GetAnilistUserOAuthToken() {
	tokenRequestUrl := fmt.Sprintf(tokenRequest, env.ANILIST_CLIENT_ID)
	fmt.Printf("Please login into your account using the following url: %s\n", tokenRequestUrl)

	openTemporaryServer()
}

func openTemporaryServer() {
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()
	srv := http.Server{
		Addr:    ":" + env.PORT,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("code")
		if len(code) == 0 {
			log.Fatal("No code was found, please try again.")
		} else {
			fmt.Fprintf(w, "You can close this page now")
			err := os.WriteFile("./.token", []byte(code), 0644)
			if err != nil {
				log.Fatal(err)
			}
			cancel()
		}
	})

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
}
