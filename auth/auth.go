package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jenrykster/gotoh/utils"
)

type CodeConversionResponse struct {
	AccessToken string `json:"access_token"`
}

const tokenRequest = "https://anilist.co/api/v2/oauth/authorize?client_id=%s&response_type=token"

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

	mux.Handle("/parseToken/", http.StripPrefix("/parseToken/", http.FileServer(http.Dir("./html/"))))

	mux.HandleFunc("/saveToken", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		code := values.Get("access_token")
		if len(code) == 0 {
			fmt.Println("There was an error while saving your token")
			log.Fatal("No code was found, please try again.")
		} else {
			fmt.Fprintf(w, "You can close this page now")
			saveAccessToken(code)
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

func saveAccessToken(code string) {
	if len(code) > 0 {
		err := os.WriteFile(utils.TOKEN_FILE_PATH, []byte(code), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success !")
	}
}
