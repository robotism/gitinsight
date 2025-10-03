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
		go AnalyzeBranchCommitLogs(insight, repoPath, branchNames)
	}
}

func AnalyzeBranchCommitLogs(insight *gitinsight.Config, repoPath string, branchNames []string) {
	repoUrl := gitinsight.GetRepoRemoteUrl(repoPath)
	repoStats, err := gitinsight.AnalyzeRepoCommitLogs(insight, repoPath, branchNames)
	if err != nil {
		log.Printf("Error analyzing repository %s: %v\n", repoPath, err)
		return
	}
	log.Printf("Analyzed %d repositories\n", len(repoStats))
	for branchName, commitLogs := range repoStats {
		log.Printf("   Analyzing repo %s branch %s\n", repoPath, branchName)
		if IsRepoUpToDate(repoPath, branchName) {
			log.Printf("   Repo %s branch %s is up to date\n", repoPath, branchName)
			continue
		}

		commitLogsModels := make([]gitinsight.CommitLogModel, len(commitLogs))
		for i, commitLog := range commitLogs {
			commitLogsModels[i] = gitinsight.CommitLogModel{
				RepoUrl:     repoUrl,
				RepoPath:    repoPath,
				BranchName:  branchName,
				CommitHash:  commitLog.Hash,
				Message:     commitLog.Message,
				MessageType: commitLog.MessageType,
				Date:        commitLog.Date,
				Additions:   commitLog.Additions,
				Deletions:   commitLog.Deletions,
				Effectives:  commitLog.Effectives,
				AuthorName:  commitLog.AuthorName,
				AuthorEmail: commitLog.AuthorEmail,
				DisplayName: commitLog.DisplayName,
			}
		}
		_, err = gitinsight.ReplaceCommitLogs(repoPath, branchName, commitLogsModels)
		if err != nil {
			log.Printf("Error clearing commit logs: %v\n", err)
			continue
		}

	}
}

func IsRepoUpToDate(repoPath string, branchName string) bool {
	log.Printf("----Checking repo %s branch %s\n", repoPath, branchName)

	localState, err := gitinsight.GetLatestCommitState(repoPath, branchName)
	if err != nil {
		log.Printf("Error getting latest commit state: %v\n", err)
		return false
	}
	log.Printf("----Local state: %v\n", localState)

	cacheCount, err := gitinsight.CountCommitLogs(&gitinsight.CommitLogFilter{
		RepoPath:   repoPath,
		BranchName: branchName,
	})
	if err != nil {
		log.Printf("Error getting latest commit state: %v\n", err)
		return false
	}
	log.Printf("----Cache count: %d\n", cacheCount)

	cacheLastestLog, err := gitinsight.GetCommitLogs(&gitinsight.CommitLogFilter{
		RepoPath:   repoPath,
		BranchName: branchName,
	}, 0, 1)
	log.Printf("----Cache latest log: %v\n", cacheLastestLog)

	if err != nil {
		log.Printf("Error getting latest commit state: %v\n", err)
		return false
	}
	if len(cacheLastestLog) != 0 {
		if localState.LatestCommitHash == cacheLastestLog[0].CommitHash && localState.CommitLogsCount == cacheCount {
			return true
		}
	}
	return false
}
