package gitinsight

import (
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
)

type CommitLogFilter struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`

	RepoUrl    string `json:"repoUrl" bun:"repo_url"`
	BranchName string `json:"branchName" bun:"branch_name"`
	CommitHash string `json:"commitHash" bun:"commit_hash"`

	IsMerge     string `json:"isMerge" bun:"is_merge"`
	MessageType string `json:"messageType" bun:"message_type"`

	DateFrom string `json:"dateFrom" bun:"date_from"`
	DateTo   string `json:"dateTo" bun:"date_to"`

	Nickname string `json:"nickname" bun:"nickname"`

	Period string `json:"period" bun:"period"`

	LeEffective string `json:"leEffective" bun:"le_effective"`
	GeEffective string `json:"geEffective" bun:"ge_effective"`
}

func (filter *CommitLogFilter) Query(query *bun.SelectQuery) {
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash = ?", filter.CommitHash)
	}
	if filter.DateFrom != "" {
		query.Where("date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query.Where("date <= ?", filter.DateTo)
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
