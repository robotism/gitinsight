package server

import (
	"log"

	"github.com/robfig/cron"
	"github.com/robotism/gitinsight/gitinsight"
)

var crond *cron.Cron
var syncing bool

func StartCrond(insight *gitinsight.Config) {
	crond = cron.New()
	crond.AddFunc("@every 1m", func() {
		if syncing {
			return
		}
		syncing = true
		defer func() {
			syncing = false
		}()
		ProcessCrond(insight)
	})
	crond.Start()
}

func StopCrond() {
	crond.Stop()
}

func ProcessCrond(insight *gitinsight.Config) {
	repos, err := gitinsight.SyncRepo(insight)
	if err != nil {
		log.Printf("Error syncing repository: %v\n", err)
		return
	}
	for repoPath, branchNames := range repos {
		repoStats, err := gitinsight.AnalyzeRepo(repoPath, branchNames)
		if err != nil {
			log.Printf("Error analyzing repository %s: %v\n", repoPath, err)
			continue
		}
		log.Printf("Synced %d repositories\n", len(repos))
		for branchName, commitLogs := range repoStats {
			log.Printf("   Synced repo %s branch %s: %v\n", repoPath, branchName, commitLogs)
		}
	}
}
