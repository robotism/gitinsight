package server

import (
	"log"

	"github.com/robotism/gitinsight/gitinsight"
)

func HandleBranchCommitLogs(insight *gitinsight.Config, repoPath string, branchNames []string) error {
	repoUrl := gitinsight.GetRepoRemoteUrl(repoPath)
	repoStats, err := gitinsight.AnalyzeRepoCommitLogs(insight, repoPath, branchNames)
	if err != nil {
		log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
		return err
	}

	for branchName, commitLogs := range repoStats {
		log.Printf("    ⏳ Caching repo %s branch %s\n", repoUrl, branchName)
		isUpToDate, err := IsRepoUpToDate(repoUrl, repoPath, branchName)
		if err != nil {
			log.Printf("❌ Error checking repo %s branch %s: %v\n", repoUrl, branchName, err)
			return err
		}
		if isUpToDate {
			log.Printf("✅   Repo %s branch %s is up to date 👍👍👍👍👍👍\n", repoUrl, branchName)
			continue
		}

		commitLogModels := make([]gitinsight.CommitLogModel, len(commitLogs))
		for i, commitLog := range commitLogs {
			commitLogModels[i] = gitinsight.CommitLogModel{
				RepoUrl:       repoUrl,
				BranchName:    branchName,
				CommitHash:    commitLog.Hash,
				IsMerge:       commitLog.IsMerge,
				Message:       commitLog.Message,
				MessageType:   commitLog.MessageType,
				Date:          commitLog.Date,
				Additions:     commitLog.Additions,
				Deletions:     commitLog.Deletions,
				Effectives:    commitLog.Effectives,
				LanguageStats: commitLog.LanguageStats,
				AuthorName:    commitLog.AuthorName,
				AuthorEmail:   commitLog.AuthorEmail,
				Nickname:      commitLog.Nickname,
			}
		}
		_, err = gitinsight.ReplaceCommitLogs(repoUrl, branchName, commitLogModels)
		if err != nil {
			log.Printf("❌ Error caching commit logs: %v\n", err)
			return err
		}
		log.Printf("✅   Cached repo %s branch commit logs\n", repoUrl)
	}
	return nil
}

func IsRepoUpToDate(repoUrl string, repoPath string, branchName string) (bool, error) {

	localState, err := gitinsight.GetLatestCommitState(repoPath, branchName)
	if err != nil {
		log.Printf("❌ Error getting latest commit state:%s %s %v\n", repoUrl, branchName, err)
		return false, err
	}
	log.Printf("    ⏳ ---- Local state:%s %s %v\n", repoUrl, branchName, localState)

	cacheCount, err := gitinsight.CountCommitLogs(&gitinsight.CommitLogFilter{
		RepoUrl:    repoUrl,
		BranchName: branchName,
	})
	if err != nil {
		log.Printf("❌  count cache commit state:%s %s %v\n", repoUrl, branchName, err)
		return false, err
	}
	log.Printf("    ⏳ ----Cache log count:%s %s %d\n", repoUrl, branchName, cacheCount)

	cacheLastestLog, err := gitinsight.GetCommitLogs(&gitinsight.CommitLogFilter{
		RepoUrl:    repoUrl,
		BranchName: branchName,
	}, 0, 1)
	log.Printf("    ⏳ ----Cache latest log:%s %s %v\n", repoUrl, branchName, cacheLastestLog)

	if err != nil {
		log.Printf("❌ Error getting cache commit state:%s %s %v\n", repoUrl, branchName, err)
		return false, err
	}
	if len(cacheLastestLog) != 0 {
		if localState.LatestCommitHash == cacheLastestLog[0].CommitHash && localState.CommitLogsCount == cacheCount {
			return true, nil
		}
	}
	return false, nil
}
