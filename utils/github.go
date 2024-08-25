package utils

import (
	"bytes"
	"code-review/config"
	"code-review/database"
	"code-review/models"
	"context"
	"fmt"
	"github.com/google/go-github/v63/github"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func prepareClientWithToken(token string) *github.Client {
	return github.NewClient(nil).WithAuthToken(token)
}

func ExchangeCodeForToken(code string) (models.GHToken, error) {
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

	var githubToken models.GHToken

	for _, kv := range strings.Split(string(responseBody), "&") {
		key := strings.Split(kv, "=")[0]
		value := strings.Split(kv, "=")[1]

		switch key {
		case "access_token":
			githubToken.AccessToken = value
		case "expires_in":
			if expiresIn, err := strconv.Atoi(value); err != nil {
				githubToken.AccessTokenExpiresIn = time.Now().Add(time.Second * time.Duration(expiresIn))
			} else {
				githubToken.AccessTokenExpiresIn = time.Now().Add(time.Hour * 8)
			}
		case "refresh_token":
			githubToken.RefreshToken = value
		case "refresh_token_expires_in":
			if expiresIn, err := strconv.Atoi(value); err != nil {
				githubToken.RefreshTokenExpiresIn = time.Now().Add(time.Second * time.Duration(expiresIn))
			} else {
				githubToken.RefreshTokenExpiresIn = time.Now().Add(time.Hour * 8)
			}
		}
	}

	return githubToken, err
}

func RefreshUserAccessToken(token string) (*models.GHToken, error) {

	var githubToken models.GHToken

	err := database.DB.Where(&githubToken, "access_token = ?", token).Error
	if err != nil {
		return nil, err
	}

	clientId := config.GetEnv("GITHUB_CLIENT_ID")
	clientSecret := config.GetEnv("GITHUB_CLIENT_SECRET")
	githubTokenURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s", clientId, clientSecret, githubToken.RefreshToken)

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

	for _, kv := range strings.Split(string(responseBody), "&") {
		key := strings.Split(kv, "=")[0]
		value := strings.Split(kv, "=")[1]

		switch key {
		case "access_token":
			githubToken.AccessToken = value
		case "expires_in":
			if expiresIn, err := strconv.Atoi(value); err != nil {
				githubToken.AccessTokenExpiresIn = time.Now().Add(time.Second * time.Duration(expiresIn))
			} else {
				githubToken.AccessTokenExpiresIn = time.Now().Add(time.Hour * 8)
			}
		case "refresh_token":
			githubToken.RefreshToken = value
		case "refresh_token_expires_in":
			if expiresIn, err := strconv.Atoi(value); err != nil {
				githubToken.RefreshTokenExpiresIn = time.Now().Add(time.Second * time.Duration(expiresIn))
			} else {
				githubToken.RefreshTokenExpiresIn = time.Now().Add(time.Hour * 8)
			}
		}
	}

	if err := database.DB.Save(&githubToken).Error; err != nil {
		return nil, err
	}

	return &githubToken, err
}

func GetUserDetails(token, username string) (github.User, error) {
	client := prepareClientWithToken(token)

	// get the public user details
	user, _, err := client.Users.Get(context.Background(), username)

	if user == nil || err != nil {
		return github.User{}, err
	}

	// if the user email isn't available publicly
	if user.Email == nil {
		emails, _, err := client.Users.ListEmails(context.Background(), nil)

		for _, email := range emails {
			// set the user email to the primary and verified email
			if email.GetVerified() && email.GetPrimary() {
				user.Email = email.Email
				break
			}
		}

		if err != nil {
			return github.User{}, err
		}
	}

	return *user, err
}

func GetUserRepos(token string) ([]*github.Repository, error) {
	client := prepareClientWithToken(token)

	var repos []*github.Repository
	var err error

	var repoOptions = new(github.RepositoryListByAuthenticatedUserOptions)
	repoOptions.Visibility = "all"
	repoOptions.Direction = "desc"
	repoOptions.Sort = "updated"

	repos, _, err = client.Repositories.ListByAuthenticatedUser(context.Background(), repoOptions)

	if err != nil {
		return repos, err
	}

	return repos, err
}
