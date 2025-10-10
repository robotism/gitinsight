package gitinsight

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/object"
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

func GetLocalCommitState(repoPath string, filter CheckUpTodateFilter) (*BranchState, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	// Try local branch first
	branchRef, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+filter.BranchName), true)
	if err != nil {
		// If local branch does not exist, try remote branch
		branchRef, err = repo.Reference(plumbing.ReferenceName("refs/remotes/origin/"+filter.BranchName), true)
		if err != nil {
			return nil, fmt.Errorf("could not get branch reference: %v", err)
		}
	}
	// Get the commit history
	cIter, err := repo.Log(&git.LogOptions{
		From: branchRef.Hash(),
	})
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("could not get commit log: %s", filter.BranchName))
	}
	count := 0
	hashcode := ""
	for {
		c, err := cIter.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if IsBeforeSince(c, filter) {
			break
		}

		isMerge := len(c.ParentHashes) > 1

		if filter.IsMerge == "1" && !isMerge {
			continue
		}
		if filter.IsMerge == "0" && isMerge {
			continue
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

func AnalyzeRepoCommitLogs(config *Config, repoPath string, filter CheckUpTodateFilter) ([]CommitLog, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	// Get branch stats
	log.Printf("🚀  Analyzing branch commit logs: %s %s\n", repoPath, filter.BranchName)
	commitLogs, err := AnalyzeBranchCommitLogs(config, repo, filter)
	if err != nil {
		log.Printf("  ⚠️ Error analyzing branch commit logs %s: %v\n", filter.BranchName, err)
		return nil, err
	}
	log.Printf("    Found %s %s %d commits\n", repoPath, filter.BranchName, len(commitLogs))
	return commitLogs, nil
}

func AnalyzeBranchCommitLogs(config *Config, repo *git.Repository, filter CheckUpTodateFilter) ([]CommitLog, error) {
	// Get the branch reference (try local first, then remote)
	var branchRef *plumbing.Reference
	var err error

	// Try local branch first
	branchRef, err = repo.Reference(plumbing.ReferenceName("refs/heads/"+filter.BranchName), true)
	if err != nil {
		// If local branch does not exist, try remote branch
		branchRef, err = repo.Reference(plumbing.ReferenceName("refs/remotes/origin/"+filter.BranchName), true)
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
		if IsBeforeSince(c, filter) {
			break
		}

		nickname := FindNickname(config, c.Author.Name, c.Author.Email)
		var additions, deletions int
		if len(c.ParentHashes) == 0 {
			// 初始提交
			additions, deletions, _ = CountLinesInCommit(c)
		} else {
			additions, deletions = GetCommitDiff(c)
		}

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
		log.Printf("    🏷️  Analyzed commit logs: %s %s %s %s %s %s %s %s\n",
			filter.RepoUrl, filter.BranchName, c.Hash.String(), nickname, c.Author.Name, c.Author.Email, c.Author.When, c.Message)
	}

	return commitLogs, nil
}

// CountLinesInCommit 统计提交中所有文件的行数
func CountLinesInCommit(commit *object.Commit) (additions int, deletions int, err error) {
	additions = 0
	deletions = 0

	tree, err := commit.Tree()
	if err != nil {
		return 0, 0, err
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		content, err := f.Contents()
		if err != nil {
			return err
		}

		// 用 bufio.Scanner 逐行统计
		scanner := bufio.NewScanner(bytes.NewReader([]byte(content)))
		count := 0
		for scanner.Scan() {
			count++
		}
		additions += count // 初始提交全部算作新增
		return scanner.Err()
	})

	return additions, deletions, err
}
