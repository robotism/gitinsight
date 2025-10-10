package server

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/robotism/gitinsight/gitinsight"
	"github.com/robotism/gitinsight/web"
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

var gConfig *AppConfig

func GetConfig() *AppConfig {
	return gConfig
}

func Run(config *AppConfig) error {
	gConfig = config
	insight := config.Insight
	server := config.Server

	err := gitinsight.OpenDb(server.Database.Type, server.Database.Dsn)
	if err != nil {
		return err
	}
	if insight.Reset {
		for i := 0; i < 3; i++ {
			log.Printf("Will reset database and repo cache, wait %d seconds...", 3-i)
			time.Sleep(1 * time.Second)
		}
		err := gitinsight.ResetDb(&insight)
		if err != nil {
			return err
		}
		err = gitinsight.ResetRepo(&insight)
		if err != nil {
			return err
		}
		log.Println("Database and repo cache reset completed.")
		time.Sleep(1 * time.Second)
	}

	err = gitinsight.InitDb()
	if err != nil {
		return err
	}

	log.Printf("load config: %v\n", config)

	if !insight.ReadOnly {
		StartCrond(&insight)
	}

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	fs := http.FS(web.WebDistFs)

	g := gin.New()

	g.Use(cors.Default())
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	res := xres.New(web.WebDistFs)

	g.NoRoute(func(c *gin.Context) {
		path := strings.TrimLeft(c.Request.URL.Path, "/")
		pathes := []string{path, "dist/" + path}
		for _, p := range pathes {
			if ok, _ := res.IsFile(p); ok {
				c.FileFromFS(p, fs)
				return
			}
		}
		c.FileFromFS("dist/", fs)
	})

	v1 := g.Group("/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	RegisterRoute(v1)

	return g.Run(server.Address)
}
