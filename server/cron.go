package server

import (
	"context"
	"log"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xgrpool"
	"github.com/robfig/cron"
	"github.com/robotism/gitinsight/gitinsight"
)

var crond *cron.Cron
var syncing bool

func StartCrond(insight *gitinsight.Config) {
	go func() {
		OnCrond(insight)
	}()
	crond = cron.New()
	crond.AddFunc("@every "+insight.Interval, func() {
		OnCrond(insight)
	})
	crond.Start()
}

func StopCrond() {
	crond.Stop()
}

func OnCrond(insight *gitinsight.Config) {
	if syncing {
		log.Printf("🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒  Sync by cron skip 🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒\n")
		return
	}
	syncing = true
	defer func() {
		syncing = false
	}()
	ProcessCrond(insight)
	log.Printf("🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒  Sync by cron done 🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒🔒\n")
	panic("sync by cron done")
}

func ProcessCrond(insight *gitinsight.Config) {
	log.Printf("⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳  Sync by cron start ⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳\n")
	timeStart := time.Now()
	repos, err := gitinsight.SyncRepo(insight)
	if err != nil {
		log.Printf("❌ Error syncing repository: %v\n", err)
		return
	}

	pool := xgrpool.New()

	log.Printf("⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳  Analyze by cron start ⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳\n")
	for repoPath, branchNames := range repos {
		for _, branchName := range branchNames {
			pool.Add(func(ctx context.Context) error {
				handleStart := time.Now()
				err := HandleBranchCommitLogs(insight, repoPath, branchName)
				if err != nil {
					log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
				}
				handleStop := time.Now()
				handleCost := handleStop.Sub(handleStart)
				log.Printf("⏰⏰⏰⏰⏰⏰  Handled %s %s cost %v ⏰⏰⏰⏰⏰⏰\n", repoPath, branchName, handleCost)
				log.Printf("✅✅✅✅✅✅  Handled %s %s done ✅✅✅✅✅✅\n", repoPath, branchName)
				return nil
			})
		}
	}
	pool.Wait()
	timeStop := time.Now()
	timeCost := timeStop.Sub(timeStart)
	log.Printf("⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰  Analyzed by cron cost %v ⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰⏰\n", timeCost)
	log.Printf("✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅  Analyzed by cron done ✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅✅\n")
}
