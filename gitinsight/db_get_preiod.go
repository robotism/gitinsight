package gitinsight

import (
	"context"
	"errors"
	"strings"

	"github.com/chaos-plus/chaos-plus-toolx/xcast"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect"
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

// GetCommitStatsByPeriodAndUser 支持按日/周/月统计，每个用户一道线，兼容 MySQL/SQLite/PostgreSQL
func GetCommitStatsByPeriodAndUser(filter *CommitLogFilter, period string) ([]CommitPeriodStatItem, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var results []CommitPeriodStatItem

	// === 根据数据库类型生成 period 表达式 ===
	var periodExpr string
	dbType := gdb.Dialect().Name()

	switch strings.ToLower(period) {
	case "day", "daily":
		switch dbType {
		case dialect.MySQL:
			periodExpr = "DATE(date)"
		case dialect.SQLite:
			periodExpr = "DATE(date)"
		case dialect.PG:
			periodExpr = "TO_CHAR(date::date, 'YYYY-MM-DD')"
		default:
			return nil, errors.New("unsupported db dialect for daily stats")
		}
	case "week", "weekly":
		switch dbType {
		case dialect.MySQL:
			periodExpr = "DATE_FORMAT(DATE_ADD(date, INTERVAL (6 - WEEKDAY(date)) DAY), '%Y-%m-%d')" // 周日
		case dialect.SQLite:
			periodExpr = "DATE(date, 'weekday 0')" // SQLite: weekday 0 = Sunday
		case dialect.PG:
			periodExpr = "TO_CHAR(date_trunc('week', date + interval '6 days')::date, 'YYYY-MM-DD')" // 周日
		default:
			return nil, errors.New("unsupported db dialect for weekly stats")
		}
	case "month", "monthly":
		switch dbType {
		case dialect.MySQL:
			periodExpr = "DATE_FORMAT(date, '%Y-%m')"
		case dialect.SQLite:
			periodExpr = "strftime('%Y-%m', date)"
		case dialect.PG:
			periodExpr = "to_char(date, 'YYYY-MM')"
		default:
			return nil, errors.New("unsupported db dialect for monthly stats")
		}
	default:
		return nil, errors.New("invalid period, must be one of: day, week, month")
	}

	// === 构建子查询，按 commit_hash 去重 ===
	subq := gdb.NewSelect().
		Model((*CommitLogModel)(nil)).
		ColumnExpr("DISTINCT commit_hash").
		Column("nickname").
		Column("date").
		Column("additions").
		Column("deletions").
		Column("effectives")

	// === 加入过滤条件 ===
	if filter.RepoUrl != "" {
		subq.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		subq.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.AuthorName != "" {
		subq.Where("author_name IN (?)", bun.In(strings.Split(filter.AuthorName, ",")))
	}
	if filter.AuthorEmail != "" {
		subq.Where("author_email IN (?)", bun.In(strings.Split(filter.AuthorEmail, ",")))
	}
	if filter.Nickname != "" {
		subq.Where("nickname IN (?)", bun.In(strings.Split(filter.Nickname, ",")))
	}
	if filter.DateFrom != "" {
		subq.Where("date >= ?", filter.DateFrom)
	}
	if filter.DateTo != "" {
		subq.Where("date <= ?", filter.DateTo)
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
		ColumnExpr("nickname").
		ColumnExpr("COUNT(commit_hash) AS commits"). // 子查询已 distinct，不需要再 DISTINCT
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives").
		ColumnExpr(periodExpr + " AS period").
		GroupExpr("period, nickname").
		OrderExpr("period ASC")

	// === 执行查询 ===
	if err := query.Scan(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}
