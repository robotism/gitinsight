package gitinsight

import (
	"strings"
	"time"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
)

type CommitLogFilter struct {
	Offset int
	Limit  int

	RepoUrl    string
	BranchName string
	CommitHash string

	IsMerge     string
	MessageType string

	SinceUTC  string
	UntilUTC  string
	SinceTime time.Time
	UntilTime time.Time

	Nickname string

	Period string

	LeEffective string
	GeEffective string
}

func (filter *CommitLogFilter) SelectQuery(query *bun.SelectQuery) {
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash = ?", filter.CommitHash)
	}
	if !filter.SinceTime.IsZero() {
		query.Where("committer_date >= ?", filter.SinceTime.Format("2006-01-02 15:04:05"))
	}
	if !filter.UntilTime.IsZero() {
		query.Where("committer_date <= ?", filter.UntilTime.Format("2006-01-02 15:04:05"))
	}
	if filter.Nickname != "" {
		query.Where("nickname IN (?)", bun.In(strings.Split(filter.Nickname, ",")))
	}
	if filter.MessageType != "" {
		query.Where("message_type IN (?)", bun.In(strings.Split(filter.MessageType, ",")))
	}
	if filter.IsMerge != "" {
		values := strings.Split(filter.IsMerge, ",")
		nums := make([]int, len(values))
		for i, v := range values {
			nums[i] = xcast.ToInt(v)
		}
		query.Where("is_merge IN (?)", bun.In(nums))
	} else {
		query.Where("is_merge = 0")
	}
	if filter.LeEffective != "" {
		query.Where("effectives <= ?", xcast.ToInt(filter.LeEffective))
	}
	if filter.GeEffective != "" {
		query.Where("effectives >= ?", xcast.ToInt(filter.GeEffective))
	}
}

func (filter *CommitLogFilter) DeleteQuery(query *bun.DeleteQuery) {
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash = ?", filter.CommitHash)
	}
	if !filter.SinceTime.IsZero() {
		query.Where("committer_date >= ?", filter.SinceTime.Format("2006-01-02 15:04:05"))
	}
	if !filter.UntilTime.IsZero() {
		query.Where("committer_date <= ?", filter.UntilTime.Format("2006-01-02 15:04:05"))
	}
	if filter.Nickname != "" {
		query.Where("nickname IN (?)", bun.In(strings.Split(filter.Nickname, ",")))
	}
	if filter.MessageType != "" {
		query.Where("message_type IN (?)", bun.In(strings.Split(filter.MessageType, ",")))
	}
	if filter.IsMerge != "" {
		values := strings.Split(filter.IsMerge, ",")
		nums := make([]int, len(values))
		for i, v := range values {
			nums[i] = xcast.ToInt(v)
		}
		query.Where("is_merge IN (?)", bun.In(nums))
	} else {
		query.Where("is_merge = 0")
	}
	if filter.LeEffective != "" {
		query.Where("effectives <= ?", xcast.ToInt(filter.LeEffective))
	}
	if filter.GeEffective != "" {
		query.Where("effectives >= ?", xcast.ToInt(filter.GeEffective))
	}
}

type CheckUpTodateFilter struct {
	RepoUrl    string
	BranchName string

	IsMerge string

	SinceUTC  string
	SinceTime time.Time
}

func (filter *CheckUpTodateFilter) ToCommitLogFilter() *CommitLogFilter {
	to := &CommitLogFilter{
		RepoUrl:    filter.RepoUrl,
		BranchName: filter.BranchName,
		IsMerge:    filter.IsMerge,
		SinceTime:  filter.SinceTime,
	}
	if !filter.SinceTime.IsZero() {
		to.SinceUTC = filter.SinceTime.UTC().Format(time.RFC3339)
	}
	return to
}
