package server

import (
	"log"

	"github.com/robfig/cron"
	"github.com/robotism/gitinsight/gitinsight"
)

var crond *cron.Cron
var syncing bool

func StartCrond(insight *gitinsight.Config) {
	go ProcessCrond(insight)
	crond = cron.New()
	crond.AddFunc("@every "+insight.Interval, func() {
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
		log.Printf("❌ Error syncing repository: %v\n", err)
		return
	}

	for repoPath, branchNames := range repos {
		if !insight.Parallel {
			err := AnalyzeBranchCommitLogs(insight, repoPath, branchNames)
			if err != nil {
				log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
			}
		} else {
			go func() {
				err := AnalyzeBranchCommitLogs(insight, repoPath, branchNames)
				if err != nil {
					log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
				}
			}()
		}
	}
}
