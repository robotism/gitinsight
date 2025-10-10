package gitinsight

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xgrpool"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

type Config struct {
	Reset    bool     `yaml:"reset" json:"reset" mapstructure:"reset" description:"reset" default:"false"`
	Parallel bool     `yaml:"parallel" json:"parallel" mapstructure:"parallel" description:"parallel" default:"true"`
	Readonly bool     `yaml:"readonly" json:"readonly" mapstructure:"readonly" description:"readonly" default:"false"`
	Interval string   `yaml:"interval" json:"interval" mapstructure:"interval" description:"interval" default:"60m"`
	Since    string   `yaml:"since" json:"since" mapstructure:"since" description:"since" default:""`
	Auths    []Auth   `yaml:"auths" json:"auths" mapstructure:"auths" description:"auths"`
	Repos    []Repo   `yaml:"repos" json:"repos" mapstructure:"repos" description:"repos"`
	Authors  []Author `yaml:"authors" json:"authors" mapstructure:"authors" description:"authors"`
	Cache    Cache    `yaml:"cache" json:"cache" mapstructure:"cache" description:"cache"`
}

func (config *Config) SinceTime() time.Time {
	if config.Since == "" {
		return time.Time{}
	}
	fmts := []string{
		time.RFC3339,
		time.RFC3339Nano,
	}
	for _, fmt := range fmts {
		t, err := time.Parse(fmt, config.Since)
		if err == nil {
			return t
		}
	}
	log.Fatalf("Invalid since time: %s", config.Since)
	return time.Time{}
}

type Auth struct {
	Domain        string `yaml:"domain" json:"domain" mapstructure:"domain" description:"domain"`
	Username      string `yaml:"username,omitempty" json:"username,omitempty" mapstructure:"username" description:"username"`
	Password      string `yaml:"password,omitempty" json:"password,omitempty" mapstructure:"password" description:"password"`
	CommitUrlTmpl string `yaml:"commit_url_tmpl,omitempty" json:"commit_url_tmpl,omitempty" mapstructure:"commit_url_tmpl" description:"commit_url_tmpl"`
}

type Repo struct {
	Url      string `yaml:"url" json:"url" mapstructure:"url" description:"url"`
	User     string `yaml:"user" json:"user" mapstructure:"user" description:"user"`
	Password string `yaml:"password" json:"password" mapstructure:"password" description:"password"`
}

type Cache struct {
	Path string `yaml:"path" json:"path" mapstructure:"path" description:"path" default:"./.repos"`
}

type Author struct {
	Name     string `yaml:"name" json:"name" mapstructure:"name" description:"name"`
	Email    string `yaml:"email" json:"email" mapstructure:"email" description:"email"`
	Nickname string `yaml:"nickname" json:"nickname" mapstructure:"nickname" description:"nickname"`
}

func ResetRepo(config *Config) error {
	err := os.RemoveAll(config.Cache.Path)
	if err != nil {
		return err
	}
	if config.Cache.Path == "" {
		config.Cache.Path = ".repos"
	}
	log.Println("Reset repository cache: " + config.Cache.Path)
	return nil
}

func SyncRepo(config *Config) (map[string][]string, error) {

	if config.Cache.Path == "" {
		config.Cache.Path = ".repos"
	}

	repoStats := make(map[string][]string)

	for i, repoInfo := range config.Repos {
		log.Printf("\n[%d/%d] Processing repository: %s\n", i+1, len(config.Repos), repoInfo.Url)

		repoName := strings.TrimSuffix(filepath.Base(repoInfo.Url), ".git")
		repoPath := filepath.Join(config.Cache.Path, repoName)

		auth, err := FindAuth(config, &repoInfo)
		if err != nil {
			return nil, err
		}
		if auth != nil {
			if repoInfo.User == "" {
				repoInfo.User = auth.Username
			}
			if repoInfo.Password == "" {
				repoInfo.Password = auth.Password
			}
		}
		// Determine which credentials to use
		username := repoInfo.User
		password := repoInfo.Password

		h := func() error {
			// Clone or update repository
			repo, err := CloneOrUpdateRepo(repoInfo.Url, repoPath, username, password)
			if err != nil {
				return fmt.Errorf("error processing %s: %v", repoInfo.Url, err)
			}
			// Get all branches
			branches, err := GetBranches(repo)
			if err != nil {
				return fmt.Errorf("error getting branches for %s: %v", repoInfo.Url, err)
			}
			log.Printf("    Found %d branches\n", len(branches))
			repoStats[repoPath] = branches
			return nil
		}
		pool := xgrpool.New()
		if config.Parallel {
			pool.AddWithRecover(func(ctx context.Context) error {
				return h()
			}, func(ctx context.Context, err interface{}) {
				log.Printf("  ⚠️ Error processing %s: %v\n", repoInfo.Url, err)
			})
		} else {
			err := h()
			if err != nil {
				return nil, fmt.Errorf("error processing %s: %v", repoInfo.Url, err)
			}
		}
		pool.Wait()

	}
	return repoStats, nil
}

func CloneOrUpdateRepo(url, path, username, password string) (*git.Repository, error) {
	// Get authentication
	var auth *http.BasicAuth

	// Priority: 1. Function parameters, 2. Environment variables
	if username != "" && password != "" {
		auth = &http.BasicAuth{
			Username: username,
			Password: password,
		}
	}

	// Check if repo exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Clone the repository with all branches
		log.Printf("Cloning %s to %s...\n", url, path)
		repo, err := git.PlainClone(path, &git.CloneOptions{
			URL:      url,
			Auth:     auth,
			Progress: os.Stdout,
		})
		if err != nil {
			return nil, err
		}

		// Fetch all remote branches
		err = repo.Fetch(&git.FetchOptions{
			RemoteName: "origin",
			Auth:       auth,
			RefSpecs:   []config.RefSpec{"refs/heads/*:refs/remotes/origin/*"},
			Progress:   os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, fmt.Errorf("could not fetch all branches: %v", err)
		}

		return repo, nil
	}

	// Open existing repository
	log.Printf("Updating %s...\n", path)
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Fetch all branches from remote
	log.Println("  Fetching all branches...")
	err = repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		RefSpecs:   []config.RefSpec{"refs/heads/*:refs/remotes/origin/*"},
		Progress:   os.Stdout,
		Force:      true,
	})

	if err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, fmt.Errorf("fetch error: %v", err)
	}

	// Try to pull current branch
	w, err := repo.Worktree()
	if err == nil {
		err = w.Pull(&git.PullOptions{
			RemoteName: "origin",
			Auth:       auth,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return nil, fmt.Errorf("pull error: %v", err)
		}
	}

	return repo, nil
}

func GetBranches(repo *git.Repository) ([]string, error) {
	var branches []string
	branchMap := make(map[string]bool)

	// Get all remote references
	remoteRefs, err := repo.References()
	if err != nil {
		return nil, err
	}

	err = remoteRefs.ForEach(func(ref *plumbing.Reference) error {
		refName := ref.Name().String()

		// Check if it's a remote branch
		if strings.HasPrefix(refName, "refs/remotes/origin/") {
			branchName := strings.TrimPrefix(refName, "refs/remotes/origin/")
			// Skip HEAD reference
			if branchName != "HEAD" {
				branchMap[branchName] = true
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Convert map to slice
	for branch := range branchMap {
		branches = append(branches, branch)
	}

	// Sort branches for consistent output
	sort.Strings(branches)

	return branches, nil
}
