package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Jenrykster/gotoh/utils"
)

type CodeConversionResponse struct {
	AccessToken string `json:"access_token"`
}

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
			convertCodeIntoAccessToken(code)
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

func convertCodeIntoAccessToken(code string) {
	url := "https://anilist.co/api/v2/oauth/token"
	body := []byte(fmt.Sprintf(`{
		"grant_type": "authorization_code",
		"client_id": "%s",
		"client_secret": "%s",
		"code": "%s",
		"redirect_uri": "http://localhost:%s"
	}`, env.ANILIST_CLIENT_ID, env.ANILIST_SECRET, code, env.PORT))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var post CodeConversionResponse
	derr := json.NewDecoder(res.Body).Decode(&post)
	if derr != nil {
		panic(derr)
	}

	if len(post.AccessToken) > 0 {
		err := os.WriteFile("./.token", []byte(post.AccessToken), 0644)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Success !")
	}
}
