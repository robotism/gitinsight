package gitinsight

import (
	"context"
	"errors"
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

	filter.Query(query)

	// ✅ 正确分组与排序（使用 Expr）
	query.GroupExpr("DATE(date)").OrderExpr("DATE(date) ASC")

	// 执行查询
	err := query.Scan(ctx, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
