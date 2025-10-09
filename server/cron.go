package server

import (
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
	if insight.Interval == "" {
		return
	}
	crond.AddFunc("@every "+insight.Interval, func() {
		OnCrond(insight)
	})
	crond.Start()
}

func StopCrond() {
	if crond == nil {
		return
	}
	crond.Stop()
}

func OnCrond(insight *gitinsight.Config) {
	if syncing {
		return
	}
	syncing = true
	defer func() {
		syncing = false
	}()
	gitinsight.HandleCommitLogs(insight)
}
