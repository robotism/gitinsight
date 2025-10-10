package gitinsight

import (
	"context"
	"log"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xgrpool"
)

func HandleCommitLogs(insight *Config) {
	log.Printf("⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳  Sync by cron start ⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳⏳\n")
	timeStart := time.Now()
	repos, err := SyncRepo(insight)
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
				err := HandleBranchCommitLogsToDb(insight, repoPath, branchName)
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

func HandleBranchCommitLogsToDb(insight *Config, repoPath string, branchName string) error {
	repoUrl := GetRepoRemoteUrl(repoPath)
	isUpToDate, err := IsRepoUpToDate(insight, repoUrl, repoPath, branchName)
	if err != nil {
		log.Printf("❌ Error checking repo %s branch %s: %v\n", repoUrl, branchName, err)
		return err
	}
	if isUpToDate {
		log.Printf("✅   Repo %s branch %s is up to date 👍👍👍👍👍👍\n", repoUrl, branchName)
		return nil
	}

	repoStats, err := AnalyzeRepoCommitLogs(insight, repoPath, []string{branchName})
	if err != nil {
		log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
		return err
	}

	for branchName, commitLogs := range repoStats {
		log.Printf("    ⏳ Caching repo %s branch %s\n", repoUrl, branchName)

		commitLogModels := make([]CommitLogModel, len(commitLogs))
		for i, commitLog := range commitLogs {
			commitLogModels[i] = CommitLogModel{
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
		_, err = ReplaceCommitLogs(repoUrl, branchName, commitLogModels)
		if err != nil {
			log.Printf("❌ Error caching commit logs: %v\n", err)
			return err
		}
		log.Printf("✅   Cached repo %s branch commit logs\n", repoUrl)
	}
	return nil
}

func IsRepoUpToDate(config *Config, repoUrl string, repoPath string, branchName string) (bool, error) {

	localState, err := GetLatestCommitState(config, repoPath, branchName)
	if err != nil {
		log.Printf("❌ Error getting latest commit state:%s %s %v\n", repoUrl, branchName, err)
		return false, err
	}
	log.Printf("    ⏳ ---- Local state:%s %s %v\n", repoUrl, branchName, localState)

	cacheCount, err := CountCommitLogs(&CommitLogFilter{
		RepoUrl:    repoUrl,
		BranchName: branchName,
		DateFrom:   config.Since,
	})
	if err != nil {
		log.Printf("❌  count cache commit state:%s %s %v\n", repoUrl, branchName, err)
		return false, err
	}
	log.Printf("    ⏳ ----Cache log count:%s %s %d\n", repoUrl, branchName, cacheCount)

	if cacheCount == 0 && localState.CommitLogsCount == 0 {
		return true, nil
	}

	cacheLastestLog, err := GetCommitLogs(&CommitLogFilter{
		Offset:     0,
		Limit:      1,
		RepoUrl:    repoUrl,
		BranchName: branchName,
		DateFrom:   config.Since,
	})
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
