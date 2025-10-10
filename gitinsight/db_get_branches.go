package gitinsight

import (
	"context"
	"errors"

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

	filter.SelectQuery(subq)

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
