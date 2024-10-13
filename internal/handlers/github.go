package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Ghat0tkach/Gohacktober-Backend/internal/github"
	"github.com/Ghat0tkach/Gohacktober-Backend/config"
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