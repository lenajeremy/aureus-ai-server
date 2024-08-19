package github

import (
	"bytes"
	"code-review/config"
	"fmt"
	"io"
	"log"
	"net/http"
)

func ExchangeCodeForToken(code string) (string, error) {
	clientId := config.GetEnv("GITHUB_CLIENT_ID")
	clientSecret := config.GetEnv("GITHUB_CLIENT_SECRET")
	githubTokenURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientId, clientSecret, code)

	reader := bytes.NewReader([]byte(""))
	resp, err := http.Post(githubTokenURL, "application/json", reader)
	if err != nil {
		log.Panic(err)
	}

	// close response body
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Panic(err)
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		log.Println("Error response. Status Code: ", resp.StatusCode)
	}

	log.Println("Response:", string(responseBody))

	return string(responseBody), err
}
