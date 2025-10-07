package gitinsight

import (
	"context"
	"errors"
	"strings"

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

func GetRepoBranches(since string, until string, repo string) ([]BranchDTO, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var branches []BranchDTO = make([]BranchDTO, 0)
	query := gdb.NewSelect().
		Model((*BranchDTO)(nil)).
		ColumnExpr("repo_url").
		ColumnExpr("branch_name").
		ColumnExpr("GROUP_CONCAT(DISTINCT nickname) AS nicknames").
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives").
		ColumnExpr("COUNT(*) AS commits")
	if since != "" {
		query.Where("date >= ?", since)
	}
	if until != "" {
		query.Where("date <= ?", until)
	}
	if repo != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(repo, ",")))
	}
	query.Group("repo_url", "branch_name")
	err := query.Distinct().Scan(ctx, &branches)
	return branches, err
}
