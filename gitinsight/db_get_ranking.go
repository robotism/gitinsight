package gitinsight

import (
	"context"
	"errors"

	"github.com/uptrace/bun"
)

type Ranking struct {
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

func GetRanking(filter *CommitLogFilter) ([]Ranking, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}

	ctx := context.Background()
	subQuery := gdb.NewSelect().
		Model((*CommitLogModel)(nil)).
		ColumnExpr("DISTINCT commit_hash, nickname, author_name, author_email, additions, deletions, effectives, repo_url, date").
		Where("is_merge = 0")

	filter.SelectQuery(subQuery)

	// 外层再统计
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

	var ranking []Ranking
	err := query.Scan(ctx, &ranking)
	return ranking, err

}
