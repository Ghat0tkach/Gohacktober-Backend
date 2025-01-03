package config

import "os"

type Config struct {
	GithubToken string
	GithubOrg   string
	FrontendURL string
}

func Load() (*Config, error) {
	return &Config{
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		GithubOrg:   os.Getenv("GITHUB_ORG"),
		FrontendURL: os.Getenv("FRONTEND_URL"),
	}, nil
}