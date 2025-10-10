package cmd

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/robotism/flagger"
	"github.com/robotism/gitinsight/gitinsight"
	"github.com/robotism/gitinsight/server"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"
)

type GenConfig struct {
	File      string `mapstructure:"file" short:"f" description:"config file" default:"config.yaml"`
	Overwrite bool   `mapstructure:"overwrite" short:"o" description:"overwrite existing config file" default:"false"`
}

var (
	genConfigFlagger = flagger.New()
	genConfig        = &GenConfig{}
)

// rootCmd represents the base command when called without any subcommands
var configGenCmd = &cobra.Command{
	Use: "gen",
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(genConfig.File); err == nil {
			if !genConfig.Overwrite {
				log.Fatalf("config file %s already exists", genConfig.File)
			}
		}
		tmpl := server.AppConfig{
			Debug: false,
			Server: server.Server{
				Address: "0.0.0.0:8080",
				Database: server.Database{
					Type: "sqliteshim",
					Dsn:  "file:gitinsight.db",
				},
			},
			Insight: gitinsight.Config{
				Interval: "15m",
				Reset:    false,
				Readonly: false,
				Parallel: true,
				Since:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.Now().Location()).Format(time.RFC3339),
				Auths: []gitinsight.Auth{
					{
						Domain:   "github.com",
						Username: "robotism",
						Password: "robotism",
					},
				},
				Repos: []gitinsight.Repo{
					{
						Url:      "https://github.com/robotism/gitinsight.git",
						User:     "robotism",
						Password: "robotism",
					},
				},
				Authors: []gitinsight.Author{
					{
						Name:     "robotism",
						Email:    "robotism@robotism.com",
						Nickname: "robotism",
					},
				},
				Cache: gitinsight.Cache{
					Path: "./.repos",
				},
			},
		}
		fmt := filepath.Ext(genConfig.File)
		if fmt == ".json" {
			json, err := json.MarshalIndent(tmpl, "", "  ")
			if err != nil {
				log.Fatalf("failed to marshal config: %v", err)
			}
			log.Printf("config file %s generated:\n%s", genConfig.File, string(json))
			if err := os.WriteFile(genConfig.File, json, 0644); err != nil {
				log.Fatalf("failed to write config file: %v", err)
			}
			return
		}
		if fmt == ".yaml" || fmt == ".yml" {
			yaml, err := yaml.Marshal(tmpl)
			if err != nil {
				log.Fatalf("failed to marshal config: %v", err)
			}
			log.Printf("config file %s generated:\n%s", genConfig.File, string(yaml))
			if err := os.WriteFile(genConfig.File, yaml, 0644); err != nil {
				log.Fatalf("failed to write config file: %v", err)
			}
			return
		}
		log.Fatalf("unsupported config file format: %s", fmt)
	},
}

// rootCmd represents the base command when called without any subcommands
var configCmd = &cobra.Command{
	Use: "config",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	genConfigFlagger.UseFlags(configGenCmd.Flags())
	genConfigFlagger.UseConfigFileArgDefault()
	genConfigFlagger.Parse(genConfig)

	configCmd.AddCommand(configGenCmd)

	rootCmd.AddCommand(configCmd)

}
