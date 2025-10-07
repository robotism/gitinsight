package gitinsight

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xgrpool"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
)

type CommitLog struct {
	Hash        string
	Message     string
	MessageType string
	IsMerge     bool
	Date        time.Time

	Additions     int
	Deletions     int
	Effectives    int
	LanguageStats string

	AuthorName  string
	AuthorEmail string
	Nickname    string
}

type BranchState struct {
	LatestCommitHash string
	CommitLogsCount  int
}

func GetLatestCommitState(config *Config, repoPath string, branchName string) (*BranchState, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	// Try local branch first
	branchRef, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+branchName), true)
	if err != nil {
		// If local branch does not exist, try remote branch
		branchRef, err = repo.Reference(plumbing.ReferenceName("refs/remotes/origin/"+branchName), true)
		if err != nil {
			return nil, fmt.Errorf("could not get branch reference: %v", err)
		}
	}
	count := 0
	// Get the commit history
	cIter, err := repo.Log(&git.LogOptions{
		From: branchRef.Hash(),
	})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("could not get commit log: %s", branchName))
	}
	hashcode := ""
	for {
		c, err := cIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if IsBeforeSince(config, c) {
			break
		}
		if hashcode == "" {
			hashcode = c.Hash.String()
		}
		count++
	}
	return &BranchState{
		LatestCommitHash: hashcode,
		CommitLogsCount:  count,
	}, nil
}

func AnalyzeRepoCommitLogs(config *Config, repoPath string, branchNames []string) (map[string][]CommitLog, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	repoStats := make(map[string][]CommitLog)
	pool := xgrpool.New()
	for _, name := range branchNames {
		pool.AddWithRecover(func(ctx context.Context) error {
			// Get branch stats
			log.Printf("🚀  Analyzing branch commit logs: %s %s\n", repoPath, name)
			commitLogs, err := AnalyzeBranchCommitLogs(config, repo, name)
			if err != nil {
				log.Printf("  ⚠️ Error analyzing branch commit logs %s: %v\n", name, err)
				return err
			}
			log.Printf("    Found %s %s %d commits\n", repoPath, name, len(commitLogs))
			repoStats[name] = commitLogs
			return nil
		}, func(ctx context.Context, err interface{}) {
			log.Printf("  ⚠️ Error analyzing branch commit logs %s: %v\n", name, err)
			panic(err)
		})
	}
	pool.Wait()
	return repoStats, nil
}

func AnalyzeBranchCommitLogs(config *Config, repo *git.Repository, branchName string) ([]CommitLog, error) {
	// Get the branch reference (try local first, then remote)
	var branchRef *plumbing.Reference
	var err error

	// Try local branch first
	branchRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+branchName), true)
	if err != nil {
		// If local branch does not exist, try remote branch
		branchRef, err = repo.Reference(plumbing.ReferenceName("refs/remotes/origin/"+branchName), true)
		if err != nil {
			return nil, fmt.Errorf("could not get branch reference: %v", err)
		}
	}

	// Get the commit history
	cIter, err := repo.Log(&git.LogOptions{
		From: branchRef.Hash(),
	})

	if err != nil {
		return nil, fmt.Errorf("could not get commit history: %v", err)
	}

	commitLogs := make([]CommitLog, 0)

	for {
		c, err := cIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if IsBeforeSince(config, c) {
			break
		}

		nickname := FindNickname(config, c.Author.Name, c.Author.Email)
		if nickname == "" {
			nickname = c.Author.Name
		}
		additions, deletions := GetCommitDiff(c)
		languageStats := GetLanguageStatPatch(c)
		languageStatsJson, _ := json.MarshalIndent(languageStats, "", "  ")
		commitLog := CommitLog{
			Hash:          c.Hash.String(),
			Message:       strings.TrimSpace(c.Message),
			MessageType:   GetMessageType(c.Message),
			IsMerge:       len(c.ParentHashes) > 1,
			Date:          c.Author.When.UTC(),
			Additions:     additions,
			Deletions:     deletions,
			Effectives:    int(math.Max(float64(additions-deletions), 0)),
			AuthorName:    c.Author.Name,
			AuthorEmail:   c.Author.Email,
			Nickname:      nickname,
			LanguageStats: string(languageStatsJson),
		}
		commitLogs = append(commitLogs, commitLog)
		log.Printf("    🏷️  Analyzed commit logs: %s %s %s %s %s %s %s\n", branchName, c.Hash.String(), nickname, c.Author.Name, c.Author.Email, c.Author.When, c.Message)
	}

	return commitLogs, nil
}
