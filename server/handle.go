package server

import (
	"log"

	"github.com/robotism/gitinsight/gitinsight"
)

func AnalyzeBranchCommitLogs(insight *gitinsight.Config, repoPath string, branchNames []string) error {
	repoUrl := gitinsight.GetRepoRemoteUrl(repoPath)
	repoStats, err := gitinsight.AnalyzeRepoCommitLogs(insight, repoPath, branchNames)
	if err != nil {
		log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
		return err
	}
	log.Printf(" Analyzed %d repositories\n", len(repoStats))
	for branchName, commitLogs := range repoStats {
		log.Printf("   Analyzing repo %s branch %s\n", repoPath, branchName)
		isUpToDate, err := IsRepoUpToDate(repoPath, branchName)
		if err != nil {
			log.Printf("Error checking repo %s branch %s: %v\n", repoPath, branchName, err)
			return err
		}
		if isUpToDate {
			log.Printf("✅   Repo %s branch %s is up to date\n", repoPath, branchName)
			continue
		}

		commitLogsModels := make([]gitinsight.CommitLogModel, len(commitLogs))
		for i, commitLog := range commitLogs {
			commitLogsModels[i] = gitinsight.CommitLogModel{
				RepoUrl:     repoUrl,
				RepoPath:    repoPath,
				BranchName:  branchName,
				CommitHash:  commitLog.Hash,
				IsMerge:     commitLog.IsMerge,
				Message:     commitLog.Message,
				MessageType: commitLog.MessageType,
				Date:        commitLog.Date,
				Additions:   commitLog.Additions,
				Deletions:   commitLog.Deletions,
				Effectives:  commitLog.Effectives,
				AuthorName:  commitLog.AuthorName,
				AuthorEmail: commitLog.AuthorEmail,
				Nickname:    commitLog.Nickname,
			}
		}
		_, err = gitinsight.ReplaceCommitLogs(repoPath, branchName, commitLogsModels)
		if err != nil {
			log.Printf("Error clearing commit logs: %v\n", err)
			return err
		}

	}
	return nil
}

func IsRepoUpToDate(repoPath string, branchName string) (bool, error) {
	log.Printf("----Checking repo %s branch %s\n", repoPath, branchName)

	localState, err := gitinsight.GetLatestCommitState(repoPath, branchName)
	if err != nil {
		log.Printf("Error getting latest commit state: %v\n", err)
		return false, err
	}
	log.Printf("----Local state: %v\n", localState)

	cacheCount, err := gitinsight.CountCommitLogs(&gitinsight.CommitLogFilter{
		RepoPath:   repoPath,
		BranchName: branchName,
	})
	if err != nil {
		log.Printf("Error count cache commit state: %v\n", err)
		return false, err
	}
	log.Printf("----Cache count: %d\n", cacheCount)

	cacheLastestLog, err := gitinsight.GetCommitLogs(&gitinsight.CommitLogFilter{
		RepoPath:   repoPath,
		BranchName: branchName,
	}, 0, 1)
	log.Printf("----Cache latest log: %v\n", cacheLastestLog)

	if err != nil {
		log.Printf("Error getting cache commit state: %v\n", err)
		return false, err
	}
	if len(cacheLastestLog) != 0 {
		if localState.LatestCommitHash == cacheLastestLog[0].CommitHash && localState.CommitLogsCount == cacheCount {
			return true, nil
		}
	}
	return false, nil
}
