package gitinsight

import (
	"context"
	"errors"
	"strings"

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
func GetCommitStatsByPeriodAndUser(filter *CommitLogFilter) ([]CommitPeriodStatItem, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var results []CommitPeriodStatItem

	// === 根据数据库类型生成 period 表达式 ===
	var periodExpr string
	dbType := gdb.Dialect().Name()

	period := filter.Period
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

	filter.Query(subq)

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
