package config

import "os"

type Config struct {
	GithubToken string
	GithubOrg   string
}

func Load() (*Config, error) {
	return &Config{
		GithubToken: os.Getenv("GITHUB_TOKEN"),
		GithubOrg:   os.Getenv("GITHUB_ORG"),
	}, nil
}