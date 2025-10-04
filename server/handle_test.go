package server

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v6"
	"github.com/goccy/go-yaml"
	"github.com/robotism/gitinsight/gitinsight"
)

func TestIsRepoUpToDate(t *testing.T) {
	workDir := "../cmd/gitinsight"
	configFile := filepath.Join(workDir, "config.yaml")
	configContent, err := os.ReadFile(configFile)
	if err != nil {
		t.Errorf("Error reading config file: %v", err)
	}
	config := AppConfig{}
	if err := yaml.Unmarshal(configContent, &config); err != nil {
		t.Errorf("Error unmarshalling config: %v", err)
	}

	fmt.Printf("load config %v\n", config)

	server := config.Server

	err = gitinsight.OpenDb(server.Database.Type, server.Database.Dsn)
	if err != nil {
		t.Errorf("Error opening db: %v", err)
	}
	err = gitinsight.InitDb()
	if err != nil {
		t.Errorf("Error init db: %v", err)
	}

	for _, repo := range config.Insight.Repos {
		repoUrl := repo.Url
		repoName := strings.TrimSuffix(filepath.Base(repoUrl), ".git")
		repoPath := filepath.Join(workDir, config.Insight.Cache.Path, repoName)

		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			t.Errorf("Error opening repo: %v", err)
		}

		branches, err := gitinsight.GetBranches(repo)
		if err != nil {
			t.Errorf("Error getting branches: %v", err)
		}

		for _, branchName := range branches {
			isUpToDate, err := IsRepoUpToDate(repoUrl, repoPath, branchName)
			if err != nil {
				t.Errorf("Error checking repo up to date: %v", err)
				panic("Error checking repo up to date")
			}
			if !isUpToDate {
				t.Errorf("Repo is not up to date: %s %s", repoUrl, branchName)
				panic("Repo is not up to date")
			} else {
				fmt.Printf("✅ branch %s %s is up to date !", repoPath, branchName)
			}
		}
	}
}
