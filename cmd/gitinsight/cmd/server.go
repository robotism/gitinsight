package cmd

import (
	"github.com/robotism/flagger"
	"github.com/robotism/gitinsight/gitinsight"
	"github.com/robotism/gitinsight/server"
	"github.com/spf13/cobra"
)

var (
	serverFlagger = flagger.New()
	serverConfig  = &server.AppConfig{
		Insight: gitinsight.Config{
			Parallel: true,
		},
	}
)

// rootCmd represents the base command when called without any subcommands
var serverCmd = &cobra.Command{
	Use: "serv",
	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		err := server.Run(serverConfig)
		panic(err)
	},
}

func init() {

	serverFlagger.UseFlags(serverCmd.Flags())
	serverFlagger.UseConfigFileArgDefault()
	serverFlagger.UseConfigPathDefault()
	serverFlagger.UseConfigTypeYaml()
	serverFlagger.Parse(serverConfig)

	rootCmd.AddCommand(serverCmd)

}
