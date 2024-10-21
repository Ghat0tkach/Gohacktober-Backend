package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"fmt"
	"context"
	"net/url"

	"github.com/Ghat0tkach/Gohacktober-Backend/internal/github"
	"github.com/Ghat0tkach/Gohacktober-Backend/config"
	"golang.org/x/oauth2"
	githubOAuth "golang.org/x/oauth2/github"
	githubAPI "github.com/google/go-github/v38/github"
)

var cfg *config.Config

func Init(c *config.Config) {
	cfg = c
}

func GetHacktoberfestContributionsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	contributions, err := github.FetchHacktoberfestContributions(cfg.GithubOrg, username)
	if err != nil {
		http.Error(w, "Failed to fetch Hacktoberfest contributions: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contributions)
}

func GitHubAuthHandler(w http.ResponseWriter, r *http.Request) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     githubOAuth.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/auth/github/callback", os.Getenv("BASE_URL")),
	}

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{"user:email"},
		Endpoint:     githubOAuth.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/auth/github/callback", os.Getenv("BASE_URL")),
	}

	code := r.URL.Query().Get("code")
	token, err := config.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange authorization code", http.StatusInternalServerError)
		return
	}

	client := githubAPI.NewClient(config.Client(r.Context(), token))
	user, _, err := client.Users.Get(r.Context(), "")
	if err != nil {
		http.Error(w, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	login := user.GetLogin()
	accessToken := token.AccessToken

	err = SaveUserInfo(r.Context(), login, accessToken)
	if err != nil {
		http.Error(w, "Failed to save user info", http.StatusInternalServerError)
		return
	}

	redirectURL, err := url.Parse("http://localhost:3000/dashboard")
	if err != nil {
		http.Error(w, "Failed to parse redirect URL", http.StatusInternalServerError)
		return
	}

	query := redirectURL.Query()
	query.Set("login", login)
	query.Set("accessToken", accessToken)
	redirectURL.RawQuery = query.Encode()

	http.Redirect(w, r, redirectURL.String(), http.StatusFound)
}

func SaveUserInfo(ctx context.Context, login, accessToken string) error {
	fmt.Printf("Saving user info: Login=%s, AccessToken=%s\n", login, accessToken)
	return nil
}