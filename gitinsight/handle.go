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
		h := func(branchName string) error {
			handleStart := time.Now()
			err := HandleBranchCommitLogsToDb(insight, repoPath, branchName)
			if err != nil {
				return err
			}
			handleStop := time.Now()
			handleCost := handleStop.Sub(handleStart)
			log.Printf("⏰⏰⏰⏰⏰⏰  Handled %s %s cost %v ⏰⏰⏰⏰⏰⏰\n", repoPath, branchName, handleCost)
			log.Printf("✅✅✅✅✅✅  Handled %s %s done ✅✅✅✅✅✅\n", repoPath, branchName)
			return nil
		}
		for _, branchName := range branchNames {
			if insight.Parallel {
				pool.Add(func(ctx context.Context) error {
					return h(branchName)
				})
			} else {
				err := h(branchName)
				if err != nil {
					log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
				}
			}
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

	filter := CheckUpTodateFilter{
		RepoUrl:    repoUrl,
		BranchName: branchName,
		SinceTime:  insight.SinceTime(),
		SinceUTC:   insight.SinceTime().Format("2006-01-02 15:04:05"),
		IsMerge:    "0",
	}
	isUpToDate, err := IsRepoUpToDate(repoPath, filter)
	if err != nil {
		log.Printf("❌ Error checking repo %s branch %s: %v\n", repoUrl, branchName, err)
		return err
	}
	if isUpToDate {
		log.Printf("✅   Repo %s branch %s is up to date 👍👍👍👍👍👍\n", repoUrl, branchName)
		return nil
	} else {
		// log.Fatal("❌   Repo branch  is not up to date ❌❌❌ \n", repoUrl, branchName)
	}

	commitLogs, err := AnalyzeRepoCommitLogs(insight, repoPath, filter)
	if err != nil {
		log.Printf("❌ Error analyzing repository %s: %v\n", repoPath, err)
		return err
	}

	log.Printf("    ⏳ Caching repo %s branch %s\n", repoUrl, branchName)

	commitLogModels := make([]CommitLogModel, len(commitLogs))
	for i, commitLog := range commitLogs {
		commitLogModels[i] = CommitLogModel{
			RepoUrl:       repoUrl,
			BranchName:    filter.BranchName,
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
	_, err = ReplaceCommitLogs(filter.ToCommitLogFilter(), commitLogModels)
	if err != nil {
		log.Printf("❌ Error caching commit logs: %v\n", err)
		return err
	}
	log.Printf("✅   Cached repo %s branch commit logs\n", repoUrl)
	return nil
}

func IsRepoUpToDate(repoPath string, filter CheckUpTodateFilter) (bool, error) {

	localState, err := GetLocalCommitState(repoPath, filter)
	if err != nil {
		log.Printf("❌ Error getting latest commit state:%s %s %v\n", repoPath, filter.BranchName, err)
		return false, err
	}
	log.Printf("    ⏳ ---- Local state:%s %s %v\n", repoPath, filter.BranchName, localState)

	cacheCount, err := CountCommitLogs(filter.ToCommitLogFilter())
	if err != nil {
		log.Printf("❌  count cache commit state:%s %s %v\n", repoPath, filter.BranchName, err)
		return false, err
	}
	log.Printf("    ⏳ ----Cache log count:%s %s %d\n", repoPath, filter.BranchName, cacheCount)

	if cacheCount == 0 && localState.CommitLogsCount == 0 {
		return true, nil
	}

	commitLogFilter := filter.ToCommitLogFilter()
	commitLogFilter.Limit = 1
	cacheLastestLog, err := GetCommitLogs(commitLogFilter)
	log.Printf("    ⏳ ----Cache latest log:%s %s %v\n", repoPath, filter.BranchName, cacheLastestLog)

	if err != nil {
		log.Printf("❌ Error getting cache commit state:%s %s %v\n", repoPath, filter.BranchName, err)
		return false, err
	}
	if len(cacheLastestLog) != 0 {
		if localState.LatestCommitHash == cacheLastestLog[0].CommitHash && localState.CommitLogsCount == cacheCount {
			return true, nil
		}
	}
	return false, nil
}
