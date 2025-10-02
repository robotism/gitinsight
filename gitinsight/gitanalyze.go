package gitinsight

import (
	"fmt"
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
	Date        time.Time
	Additions   int
	Deletions   int
	AuthorName  string
	AuthorEmail string
}

func AnalyzeRepo(repoPath string, branchNames []string) (map[string][]CommitLog, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	repoStats := make(map[string][]CommitLog)

	for j, branch := range branchNames {
		fmt.Printf("  [%d/%d] Analyzing branch: %s\n", j+1, len(branchNames), branch)
		// Get branch stats
		commitLogs, err := AnalyzeBranch(repo, branch)
		if err != nil {
			log.Printf("  ⚠ Error analyzing branch %s: %v\n", branch, err)
			continue
		}
		fmt.Printf("  ✓ Found %d commits\n", len(commitLogs))
		repoStats[branch] = commitLogs
	}

	return repoStats, nil
}

func AnalyzeBranch(repo *git.Repository, branchName string) ([]CommitLog, error) {
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

	// Process each commit
	err = cIter.ForEach(func(c *object.Commit) error {

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
			Date:        c.Author.When,
			Additions:   additions,
			Deletions:   deletions,
			AuthorName:  c.Author.Name,
			AuthorEmail: c.Author.Email,
		}
		commitLogs = append(commitLogs, commitLog)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error processing commits: %v", err)
	}

	return commitLogs, nil
}
