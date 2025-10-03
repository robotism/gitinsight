package gitinsight

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type CommitLog struct {
	Hash        string
	Message     string
	MessageType string
	Date        time.Time
	Additions   int
	Deletions   int
	Effectives  int
	AuthorName  string
	AuthorEmail string
	DisplayName string
}

type BranchState struct {
	LatestCommitHash string
	CommitLogsCount  int
}

func GetLatestCommitState(repoPath string, branchName string) (*BranchState, error) {
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
	hash, err := repo.ResolveRevision(plumbing.Revision(branchRef.Name().String()))
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("could not resolve revision: %s", branchName))
	}
	commit, err := repo.CommitObject(*hash)
	if err != nil {
		return nil, errors.Join(err, fmt.Errorf("could not get commit object: %s", branchName))
	}
	count := 0
	cIter, err := repo.Log(&git.LogOptions{From: commit.Hash})
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

	for _, name := range branchNames {
		fmt.Printf("🚀  Analyzing branch: %s\n", name)
		// Get branch stats
		commitLogs, err := AnalyzeBranchCommitLogs(config, repo, name)
		if err != nil {
			log.Printf("  ⚠️ Error analyzing branch %s: %v\n", name, err)
			continue
		}
		fmt.Printf("  ✓ Found %d commits\n", len(commitLogs))
		repoStats[name] = commitLogs
	}

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
	cIter, err := repo.Log(&git.LogOptions{From: branchRef.Hash()})
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
			return commitLogs, err
		}

		// Get diff stats
		var fileStats object.FileStats
		if c.NumParents() > 0 {
			parent, err := c.Parents().Next()
			if err == nil {
				parentTree, _ := parent.Tree()
				commitTree, _ := c.Tree()
				changes, _ := object.DiffTree(parentTree, commitTree)
				patch, _ := changes.Patch()
				if patch != nil {
					fileStats = patch.Stats()
				}
			}
		}

		additions, deletions := 0, 0
		for _, stat := range fileStats {
			additions += stat.Addition
			deletions += stat.Deletion
		}

		// Update commit stats
		commitLog := CommitLog{
			Hash:        c.Hash.String(),
			Message:     strings.TrimSpace(c.Message),
			MessageType: GetMessageType(c.Message),
			Date:        c.Author.When,
			Additions:   additions,
			Deletions:   deletions,
			Effectives:  additions - deletions,
			AuthorName:  c.Author.Name,
			AuthorEmail: c.Author.Email,
			DisplayName: FindNickname(config, c.Author.Name, c.Author.Email),
		}
		commitLogs = append(commitLogs, commitLog)
	}

	return commitLogs, nil
}
