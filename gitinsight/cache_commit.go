package gitinsight

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"

	_ "github.com/uptrace/bun/driver/sqliteshim"

	_ "github.com/go-sql-driver/mysql"
)

// User model
type CommitLogModel struct {
	bun.BaseModel `bun:"table:commit_logs,alias:cl"`

	ID      int64  `bun:"id,pk,autoincrement"`
	RepoUrl string `bun:",notnull"`

	BranchName  string `bun:",notnull"`
	CommitHash  string `bun:",notnull"`
	IsMerge     bool   `bun:",notnull"`
	Message     string `bun:",notnull,type:text"`
	MessageType string `bun:",notnull"`

	Date time.Time `bun:",notnull"`

	Additions     int    `bun:",notnull"`
	Deletions     int    `bun:",notnull"`
	Effectives    int    `bun:",notnull"`
	LanguageStats string `bun:",notnull,type:text"`

	AuthorName  string `bun:",notnull"`
	AuthorEmail string `bun:",notnull"`
	Nickname    string `bun:",notnull"`
}

type CommitLogFilter struct {
	RepoUrl    string `bun:"repo_url"`
	BranchName string `bun:"branch_name"`
	CommitHash string `bun:"commit_hash"`

	DateFrom time.Time `bun:"date_from"`
	DateTo   time.Time `bun:"date_to"`

	AuthorName  string `bun:"author_name"`
	AuthorEmail string `bun:"author_email"`
	Nickname    string `bun:"nickname"`
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
	for i := 0; i < len(commitLogs); i += commitLogLimit {
		end := i + commitLogLimit
		if end > len(commitLogs) {
			end = len(commitLogs)
		}
		segment := commitLogs[i:end]
		result, err := gdb.NewInsert().Model(&segment).Exec(ctx)
		if err != nil {
			return 0, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		rowsAffected += rows
	}
	return rowsAffected, nil
}

func CountCommitLogs(filter *CommitLogFilter) (int, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	ctx := context.Background()
	query := gdb.NewSelect().Model(&CommitLogModel{})
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash = ?", filter.CommitHash)
	}
	if !filter.DateFrom.IsZero() {
		query.Where("date >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query.Where("date <= ?", filter.DateTo)
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
	return query.Count(ctx)
}

func GetCommitLogs(filter *CommitLogFilter, offset int, limit int) ([]CommitLogModel, error) {
	if gdb == nil {
		return nil, errors.New("database not initialized")
	}
	ctx := context.Background()
	var commitLogs []CommitLogModel = make([]CommitLogModel, 0)
	query := gdb.NewSelect().Model(&CommitLogModel{})
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash IN (?)", bun.In(strings.Split(filter.CommitHash, ",")))
	}
	if !filter.DateFrom.IsZero() {
		query.Where("date >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query.Where("date <= ?", filter.DateTo)
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
	query.Order("date DESC").Offset(offset).Limit(limit)
	err := query.Scan(ctx, &commitLogs)
	return commitLogs, err
}

func DelCommitLog(filter *CommitLogFilter) (int64, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	ctx := context.Background()
	query := gdb.NewDelete().Model(&CommitLogModel{})
	if filter.RepoUrl != "" {
		query.Where("repo_url IN (?)", bun.In(strings.Split(filter.RepoUrl, ",")))
	}
	if filter.BranchName != "" {
		query.Where("branch_name IN (?)", bun.In(strings.Split(filter.BranchName, ",")))
	}
	if filter.CommitHash != "" {
		query.Where("commit_hash IN (?)", bun.In(strings.Split(filter.CommitHash, ",")))
	}
	if !filter.DateFrom.IsZero() {
		query.Where("date >= ?", filter.DateFrom)
	}
	if !filter.DateTo.IsZero() {
		query.Where("date <= ?", filter.DateTo)
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
	result, err := query.Exec(ctx)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
