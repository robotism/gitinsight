package gitinsight

import (
	"context"
	"errors"
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
)

type CommitHeatmapItem struct {
	Date       string `bun:"date" json:"date"`
	Commits    int    `bun:"commits" json:"commits"`
	Additions  int    `bun:"additions" json:"additions"`
	Deletions  int    `bun:"deletions" json:"deletions"`
	Effectives int    `bun:"effectives" json:"effectives"`
}

func GetCommitHeatmapData(filter *CommitLogFilter) ([]CommitHeatmapItem, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}

	ctx := context.Background()
	var results []CommitHeatmapItem

	query := gdb.NewSelect().
		Model((*CommitLogModel)(nil)).
		ColumnExpr("DATE(date) AS date").
		ColumnExpr("COUNT(DISTINCT commit_hash) AS commits").
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives")

	// 🔹 动态过滤条件
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.AuthorName != "" {
		query.Where("author_name IN (?)", bun.In(strings.Split(filter.AuthorName, ",")))
	}
	if filter.AuthorEmail != "" {
		query.Where("author_email IN (?)", bun.In(strings.Split(filter.AuthorEmail, ",")))
	}
	if filter.Nickname != "" {
		query.Where("nickname IN (?)", bun.In(strings.Split(filter.Nickname, ",")))
	}
	if filter.DateFrom != "" {
		query.Where("date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		query.Where("date <= ?", filter.DateTo)
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

	// ✅ 正确分组与排序（使用 Expr）
	query.GroupExpr("DATE(date)").OrderExpr("DATE(date) ASC")

	// 执行查询
	err := query.Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
