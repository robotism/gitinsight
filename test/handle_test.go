package gitinsight_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chaos-plus/chaos-plus-toolx/xgrpool"
	"github.com/go-git/go-git/v6"
	"github.com/goccy/go-yaml"
	"github.com/robotism/gitinsight/gitinsight"
	"github.com/robotism/gitinsight/server"
	"github.com/stretchr/testify/require"
)

func TestIsRepoUpToDate(t *testing.T) {
	t.Parallel() // 并发安全起点
	workDir := "../cmd/gitinsight"
	configFile := filepath.Join(workDir, "config.yaml")

	configContent, err := os.ReadFile(configFile)
	require.NoError(t, err, "read config failed")

	var config server.AppConfig
	require.NoError(t, yaml.Unmarshal(configContent, &config), "unmarshal failed")

	require.NoError(t, gitinsight.OpenDb(config.Server.Database.Type, config.Server.Database.Dsn))
	require.NoError(t, gitinsight.InitDb())

	pool := xgrpool.New()
	errCh := make(chan error, 100)

	for _, repoCfg := range config.Insight.Repos {
		repoUrl := repoCfg.Url
		repoName := strings.TrimSuffix(filepath.Base(repoUrl), ".git")
		repoPath := filepath.Join(workDir, config.Insight.Cache.Path, repoName)

		repo, err := git.PlainOpen(repoPath)
		require.NoError(t, err)

		branches, err := gitinsight.GetBranches(repo)
		require.NoError(t, err)

		for _, branch := range branches {
			branch := branch // 避免闭包坑
			pool.Add(func(ctx context.Context) error {
				isUp, err := gitinsight.IsRepoUpToDate(&config.Insight, repoUrl, repoPath, branch)
				if err != nil {
					errCh <- fmt.Errorf("check failed: %s %s: %v", repoUrl, branch, err)
					return nil
				}
				if !isUp {
					errCh <- fmt.Errorf("repo not up to date: %s %s", repoUrl, branch)
				} else {
					fmt.Printf("✅ %s %s up to date!\n", repoPath, branch)
				}
				return nil
			})
		}
	}

	pool.Wait()
	close(errCh)

	for err := range errCh {
		t.Error(err)
	}
}
