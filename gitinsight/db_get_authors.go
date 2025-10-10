package gitinsight

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

type AuthorDTO struct {
	bun.BaseModel `bun:"table:commit_log,alias:cl"`

	Name     string `json:"name" bun:",notnull"`
	Email    string `json:"email" bun:",notnull"`
	Nickname string `json:"nickname" bun:",notnull"`

	Additions  int `json:"additions" bun:",notnull"`
	Deletions  int `json:"deletions" bun:",notnull"`
	Effectives int `json:"effectives" bun:",notnull"`

	Projects int `json:"projects" bun:",notnull"`
	Commits  int `json:"commits" bun:",notnull"`
}

func GetAuthors(filter *CommitLogFilter) ([]AuthorDTO, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}

	ctx := context.Background()
	var authors []AuthorDTO

	// 先构建子查询，去重 commit_hash
	subQuery := gdb.NewSelect().
		Model((*CommitLogModel)(nil)).
		ColumnExpr("DISTINCT commit_hash, nickname, author_name, author_email, additions, deletions, effectives, repo_url, date")

	filter.SelectQuery(subQuery)

	// 外层统计
	query := gdb.NewSelect().
		TableExpr("(?) AS t", subQuery).
		ColumnExpr("nickname").
		ColumnExpr("GROUP_CONCAT(DISTINCT author_name) AS name").
		ColumnExpr("GROUP_CONCAT(DISTINCT author_email) AS email").
		ColumnExpr("SUM(additions) AS additions").
		ColumnExpr("SUM(deletions) AS deletions").
		ColumnExpr("SUM(effectives) AS effectives").
		ColumnExpr("COUNT(DISTINCT repo_url) AS projects").
		ColumnExpr("COUNT(DISTINCT commit_hash) AS commits").
		Group("nickname")

	err := query.Scan(ctx, &authors)
	return authors, err
}
