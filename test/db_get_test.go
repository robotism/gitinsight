package gitinsight_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/robotism/gitinsight/gitinsight"
	"github.com/robotism/gitinsight/server"
	"github.com/stretchr/testify/require"
)

func TestGetRepoBranches(t *testing.T) {

	workDir := "../cmd/gitinsight"
	configFile := filepath.Join(workDir, "config.yaml")

	configContent, err := os.ReadFile(configFile)
	require.NoError(t, err, "read config failed")

	var config server.AppConfig
	require.NoError(t, yaml.Unmarshal(configContent, &config), "unmarshal failed")

	require.NoError(t, gitinsight.OpenDb(config.Server.Database.Type, config.Server.Database.Dsn))
	require.NoError(t, gitinsight.InitDb())

	since, _ := time.Parse("2006-01-02 15:04:05", "2025-09-01 00:00:00")
	until, _ := time.Parse("2006-01-02 15:04:05", "2025-09-30 23:59:59")
	branches, err := gitinsight.GetRepoBranches(&gitinsight.CommitLogFilter{
		SinceTime: since,
		UntilTime: until,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(branches) == 0 {
		t.Fatal("no branches")
	}
	t.Log(branches)
}
