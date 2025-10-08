package gitinsight

import (
	"context"
	"errors"
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
)

// CommitPeriodStatItem 表示按天/周/月统计的数据
type CommitPeriodStatItem struct {
	Period     string `bun:"period" json:"period"`       // 日期/周/月标识
	Nickname   string `bun:"nickname" json:"nickname"`   // 提交人
	Commits    int    `bun:"commits" json:"commits"`     // 提交次数
	Additions  int    `bun:"additions" json:"additions"` // 新增行
	Deletions  int    `bun:"deletions" json:"deletions"` // 删除行
	Effectives int    `bun:"effectives" json:"effectives"`
}

// GetCommitStatsByPeriodAndUser 支持按日/周/月统计，每个用户一道线
func GetCommitStatsByPeriodAndUser(filter *CommitLogFilter, period string) ([]CommitPeriodStatItem, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var results []CommitPeriodStatItem

	query := gdb.NewSelect().Model((*CommitLogModel)(nil)).
		ColumnExpr("nickname").
		ColumnExpr("COUNT(DISTINCT commit_hash) AS commits").
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives")

	// 横坐标按时间
	switch strings.ToLower(period) {
	case "day", "daily":
		// 按天显示
		query.ColumnExpr("DATE(date) AS period").
			GroupExpr("DATE(date), nickname").
			OrderExpr("DATE(date) ASC")
	case "week", "weekly":
		// 使用周一作为每周的 period
		query.ColumnExpr("DATE(date, 'weekday 1', '-6 days') AS period").
			GroupExpr("DATE(date, 'weekday 1', '-6 days'), nickname").
			OrderExpr("DATE(date, 'weekday 1', '-6 days') ASC")
	case "month", "monthly":
		// 按月显示
		query.ColumnExpr("strftime('%Y-%m', date) AS period").
			GroupExpr("strftime('%Y-%m', date), nickname").
			OrderExpr("strftime('%Y-%m', date) ASC")
	default:
		return nil, errors.New("invalid period, must be one of: day, week, month")
	}

	// 过滤条件
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

	if err := query.Scan(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
