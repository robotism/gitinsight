package gitinsight

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/uptrace/bun"

	_ "github.com/uptrace/bun/driver/sqliteshim"

	_ "github.com/go-sql-driver/mysql"
)

// User model
type CommitLogModel struct {
	bun.BaseModel `bun:"table:commit_log,alias:cl"`

	ID      int64  `json:"id" bun:"id,pk,autoincrement"`
	RepoUrl string `json:"repoUrl" bun:",notnull"`

	BranchName  string `json:"branchName" bun:",notnull"`
	CommitHash  string `json:"commitHash" bun:",notnull"`
	IsMerge     bool   `json:"isMerge" bun:",notnull"`
	Message     string `json:"message" bun:",notnull,type:text"`
	MessageType string `json:"messageType" bun:",notnull"`

	Date time.Time `json:"date" bun:",notnull"`

	Additions     int    `json:"additions" bun:",notnull"`
	Deletions     int    `json:"deletions" bun:",notnull"`
	Effectives    int    `json:"effectives" bun:",notnull"`
	LanguageStats string `json:"languageStats" bun:",notnull,type:text"`

	AuthorName  string `json:"authorName" bun:",notnull"`
	AuthorEmail string `json:"authorEmail" bun:",notnull"`
	Nickname    string `json:"nickname" bun:",notnull"`
}

func InitCommit() error {
	ctx := context.Background()
	// Create table
	_, err := gdb.NewCreateTable().Model((*CommitLogModel)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	indexes := map[string]string{
		"idx_repo_url":     "repo_url",
		"idx_branch_name":  "branch_name",
		"idx_commit_hash":  "commit_hash",
		"idx_date":         "date",
		"idx_author_name":  "author_name",
		"idx_author_email": "author_email",
		"idx_nickname":     "nickname",
	}
	// Create indexes
	for indexName, columnName := range indexes {
		_, err = gdb.NewCreateIndex().Model((*CommitLogModel)(nil)).Index(indexName).Column(columnName).IfNotExists().Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReplaceCommitLogs(repoUrl string, branchName string, commitLogs []CommitLogModel) (int64, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	ctx := context.Background()
	var rowsAffected int64
	err := gdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().Model(&CommitLogModel{}).Where("repo_url = ? AND branch_name = ?", repoUrl, branchName).Exec(ctx)
		if err != nil {
			return err
		}

		const commitLogLimit = 1000

		var rowsAffected int64
		for i := 0; i < len(commitLogs); i += commitLogLimit {
			end := i + commitLogLimit
			if end > len(commitLogs) {
				end = len(commitLogs)
			}
			segment := commitLogs[i:end]
			result, err := tx.NewInsert().Model(&segment).Exec(ctx)
			if err != nil {
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}
			rowsAffected += rows
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func AddCommitLogs(commitLogs []CommitLogModel) (int64, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	if len(commitLogs) == 0 {
		return 0, nil
	}
	ctx := context.Background()

	const commitLogLimit = 1000

	var rowsAffected int64
	err := gdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		for i := 0; i < len(commitLogs); i += commitLogLimit {
			end := i + commitLogLimit
			if end > len(commitLogs) {
				end = len(commitLogs)
			}
			segment := commitLogs[i:end]
			result, err := tx.NewInsert().Model(&segment).Exec(ctx)
			if err != nil {
				return err
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return err
			}
			rowsAffected += rows
		}
		return nil
	})
	return rowsAffected, err
}

func CountCommitLogs(filter *CommitLogFilter) (int, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	ctx := context.Background()
	query := gdb.NewSelect().Model(&CommitLogModel{})
	filter.Query(query)
	return query.Count(ctx)
}

func GetCommitLogs(filter *CommitLogFilter) ([]CommitLogModel, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var commitLogs []CommitLogModel = make([]CommitLogModel, 0)
	query := gdb.NewSelect().Model(&CommitLogModel{})
	filter.Query(query)
	query.Order("date DESC").Offset(filter.Offset).Limit(filter.Limit)
	err := query.Scan(ctx, &commitLogs)
	return commitLogs, err
}

func ResetCommit() error {
	if gdb == nil {
		return errors.New("database not initialized")
	}
	ctx := context.Background()
	_, err := gdb.NewDropTable().IfExists().Model(&CommitLogModel{}).Exec(ctx)
	if err != nil {
		return err
	}
	log.Println("Reset commit log")
	return nil
}
