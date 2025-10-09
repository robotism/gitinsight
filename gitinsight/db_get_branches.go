package gitinsight

import (
	"context"
	"errors"
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
)

type BranchDTO struct {
	bun.BaseModel `bun:"table:commit_log,alias:cl"`

	RepoUrl    string `json:"repoUrl" bun:",notnull"`
	BranchName string `json:"branchName" bun:",notnull"`
	Commits    int    `json:"commits"`
	Nicknames  string `json:"nicknames"`
	Additions  int    `json:"additions"`
	Deletions  int    `json:"deletions"`
	Effectives int    `json:"effectives"`
}

func GetRepoBranches(filter *CommitLogFilter) ([]BranchDTO, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}

	ctx := context.Background()
	var branches []BranchDTO

	subq := gdb.NewSelect().
	Model((*CommitLogModel)(nil)).
	ColumnExpr("DISTINCT repo_url, branch_name, commit_hash").
	Column("nickname").
	Column("additions").
	Column("deletions").
	Column("effectives")


	// === 过滤条件 ===
	if filter.DateFrom != "" {
		subq.Where("date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		subq.Where("date <= ?", filter.DateTo)
	}
	if filter.RepoUrl != "" {
		subq.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		subq.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.Nickname != "" {
		subq.Where("nickname IN (?)", bun.In(strings.Split(filter.Nickname, ",")))
	}
	if filter.MessageType != "" {
		subq.Where("message_type IN (?)", bun.In(strings.Split(filter.MessageType, ",")))
	}
	if filter.IsMerge != "" {
		values := strings.Split(filter.IsMerge, ",")
		nums := make([]int, len(values))
		for i, v := range values {
			nums[i] = xcast.ToInt(v)
		}
		subq.Where("is_merge IN (?)", bun.In(nums))
	} else {
		subq.Where("is_merge = 0")
	}

	// === 外层统计 ===
	query := gdb.NewSelect().
		TableExpr("(?) AS t", subq).
		ColumnExpr("repo_url").
		ColumnExpr("branch_name").
		ColumnExpr("GROUP_CONCAT(DISTINCT nickname) AS nicknames").
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives").
		ColumnExpr("COUNT(commit_hash) AS commits").
		Group("repo_url", "branch_name")

	err := query.Scan(ctx, &branches)
	return branches, err
}
