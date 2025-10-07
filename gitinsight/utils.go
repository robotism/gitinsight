package gitinsight

import (
	"io"
	"net/url"
	"path/filepath"
	"strings"

	"time"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/object"
)

func IsBeforeSince(config *Config, c *object.Commit) bool {
	if config.Since != "" {
		since, err := time.Parse("2006-01-02 15:04:05", config.Since)
		if err == nil {
			if c.Author.When.Before(since) {
				return true
			}
		}
	}
	return false
}

func FindAuth(config *Config, repo *Repo) (*Auth, error) {
	for _, auth := range config.Auths {
		uri, err := url.Parse(repo.Url)
		if err != nil {
			return nil, err
		}
		host := strings.ToLower(uri.Host)
		domain := strings.ToLower(auth.Domain)
		if strings.Contains(host, domain) {
			return &auth, nil
		}
	}
	return nil, nil
}

func FindNickname(config *Config, authorName string, authorEmail string) string {
	for _, author := range config.Authors {
		if author.Name == authorName || author.Email == authorEmail {
			return author.Nickname
		}
	}
	return ""
}

func GetRepoRemoteUrl(repoPath string) string {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return ""
	}
	remote, err := repo.Remote("origin")
	if err != nil {
		return ""
	}
	return remote.Config().URLs[0]
}

func GetMessageType(message string) string {
	message = strings.ReplaceAll(message, "：", ":")
	if !strings.Contains(message, ":") {
		return ""
	}
	slits := strings.Split(message, ":")
	spec := strings.TrimSpace(slits[0])

	return extractLetters(spec)
}

func IsAsciiLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

func extractLetters(s string) string {
	var result []rune
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			result = append(result, r)
		} else {
			break
		}
	}
	return string(result)
}
func GetLanguageStatPatch(c *object.Commit) map[string]int {
	languageStats := make(map[string]int)

	// 如果有 parent，拿 diff
	if c.NumParents() > 0 {
		parent, err := c.Parent(0)
		if err != nil {
			return languageStats
		}

		patch, err := parent.Patch(c)
		if err != nil {
			return languageStats
		}

		for _, fileStat := range patch.FilePatches() {
			from, to := fileStat.Files()
			var filename string
			if to != nil {
				filename = to.Path()
			} else if from != nil {
				filename = from.Path()
			}

			ext := filepath.Ext(filename)
			languageStats[ext]++
		}
	} else {
		// 第一个 commit，没有 parent，就遍历所有文件
		fIter, _ := c.Files()
		fIter.ForEach(func(f *object.File) error {
			ext := filepath.Ext(f.Name)
			languageStats[ext]++
			return nil
		})
	}

	return languageStats
}

func GetLanguageStatsALl(c *object.Commit) map[string]int {
	languageStats := make(map[string]int)
	f, err := c.Files()
	if err != nil {
		return languageStats
	}
	f.ForEach(func(f *object.File) error {
		if f == nil {
			return nil
		}
		languageStats[filepath.Ext(f.Name)]++
		return nil
	})
	return languageStats
}

func GetCommitDiff(c *object.Commit) (int, int) {

	// Get diff stats
	additions, deletions := 0, 0
	// var fileStats object.FileStats
	// if c.NumParents() > 0 {
	// 	parent, err := c.Parents().Next()
	// 	if err == nil {
	// 		parentTree, _ := parent.Tree()
	// 		commitTree, _ := c.Tree()
	// 		changes, _ := object.DiffTree(parentTree, commitTree)
	// 		patch, _ := changes.Patch()
	// 		if patch != nil {
	// 			fileStats = patch.Stats()
	// 		}
	// 	}
	// }
	// for _, stat := range fileStats {
	// 	additions += stat.Addition
	// 	deletions += stat.Deletion
	// }
	if c.NumParents() > 0 {
		parentIter := c.Parents()
		for {
			parent, err := parentIter.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}
			parentTree, _ := parent.Tree()
			commitTree, _ := c.Tree()
			changes, _ := object.DiffTree(parentTree, commitTree)
			patch, _ := changes.Patch()
			if patch != nil {
				stats := patch.Stats()
				for _, stat := range stats {
					additions += stat.Addition
					deletions += stat.Deletion
				}
			}
		}
	}
	return additions, deletions
}
