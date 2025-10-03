package gitinsight

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func FindAuth(config *Config, repo *Repo) *Auth {
	var first *Auth
	for _, auth := range config.Auths {
		if first == nil {
			first = &auth
		}
		host := strings.ToLower(strings.Split(repo.Url, "://")[0])
		domain := strings.ToLower(auth.Domain)
		if strings.Contains(host, domain) {
			return &auth
		}
	}
	return first
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
