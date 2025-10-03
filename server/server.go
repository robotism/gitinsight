package server

import (
	"github.com/gin-gonic/gin"
	"github.com/robotism/gitinsight/gitinsight"
)

type AppConfig struct {
	Debug   bool              `mapstructure:"debug" short:"d" description:"debug mode" default:"false"`
	Server  Server            `mapstructure:"server" group:"server"`
	Insight gitinsight.Config `mapstructure:"insight" group:"insight"`
}

type Server struct {
	Address string `mapstructure:"address" description:"address" default:"0.0.0.0:8080"`

	Database Database `mapstructure:"database" group:"database" `
}

type Database struct {
	Type string `mapstructure:"type" description:"database type" default:"sqliteshim"`
	Dsn  string `mapstructure:"dsn" description:"database dsn" default:"./gitinsight.db"`
}

func Run(config *AppConfig) error {
	insight := config.Insight
	server := config.Server

	err := gitinsight.OpenDb(server.Database.Type, server.Database.Dsn)
	if err != nil {
		return err
	}
	err = gitinsight.InitDb()
	if err != nil {
		return err
	}

	StartCrond(&insight)

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	g.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World")
	})
	g.Run(server.Address)

	return nil
}
