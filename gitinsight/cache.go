package gitinsight

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

// User model
type CommitLogModel struct {
	bun.BaseModel `bun:"table:commit_logs,alias:cl"`

	ID          int64     `bun:"id,pk,autoincrement"`
	RepoUrl     string    `bun:",notnull"`
	RepoPath    string    `bun:",notnull"`
	BranchName  string    `bun:",notnull"`
	CommitHash  string    `bun:",notnull"`
	Message     string    `bun:",notnull"`
	MessageType string    `bun:",notnull"`
	Date        time.Time `bun:",notnull"`
	Additions   int       `bun:",notnull"`
	Deletions   int       `bun:",notnull"`
	Effectives  int       `bun:",notnull"`
	AuthorName  string    `bun:",notnull"`
	AuthorEmail string    `bun:",notnull"`
	DisplayName string    `bun:",notnull"`
}

type CommitLogFilter struct {
	RepoUrl    string `bun:"repo_url"`
	RepoPath   string `bun:"repo_path"`
	BranchName string `bun:"branch_name"`
	CommitHash string `bun:"commit_hash"`

	DateFrom time.Time `bun:"date_from"`
	DateTo   time.Time `bun:"date_to"`

	AuthorName  string `bun:"author_name"`
	AuthorEmail string `bun:"author_email"`
	Nickname    string `bun:"nickname"`
}

var gdb *bun.DB

func OpenDb() error {
	ctx := context.Background()

	// Open database connection
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:gitinsight.db")
	if err != nil {
		return err
	}

	// Create Bun database instance
	db := bun.NewDB(sqldb, sqlitedialect.New())

	// Add query debugging (optional)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))

	// Create table
	_, err = db.NewCreateTable().Model((*CommitLogModel)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}

	// Create indexes
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_repo_url").Column("repo_url").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_repo_path").Column("repo_path").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_branch_name").Column("branch_name").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_commit_hash").Column("commit_hash").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_date").Column("date").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_author_name").Column("author_name").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_author_email").Column("author_email").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	_, err = db.NewCreateIndex().Model((*CommitLogModel)(nil)).Index("idx_commit_logs_nickname").Column("nickname").IfNotExists().Exec(ctx)
	if err != nil {
		return err
	}
	gdb = db
	return nil
}

func CloseDb() error {
	if gdb == nil {
		return nil
	}
	return gdb.Close()
}

func ReplaceCommitLogs(repoPath string, branchName string, commitLogs []CommitLogModel) (int64, error) {
	if gdb == nil {
		return 0, errors.New("database not initialized")
	}
	ctx := context.Background()
	var rowsAffected int64
	gdb.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewDelete().Model(&CommitLogModel{}).Where("repo_path = ? AND branch_name = ?", repoPath, branchName).Exec(ctx)
		if err != nil {
			return err
		}
		result, err := tx.NewInsert().Model(&commitLogs).Exec(ctx)
		if err != nil {
			return err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return err
		}
		rowsAffected += rows
		return nil
	})
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

	result, err := gdb.NewInsert().Model(&commitLogs).Exec(ctx)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, nil
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
	if filter.RepoPath != "" {
		query.Where("repo_path IN (?)", bun.In(strings.Split(filter.RepoPath, ",")))
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
	if filter.RepoPath != "" {
		query.Where("repo_path IN (?)", bun.In(strings.Split(filter.RepoPath, ",")))
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
	if filter.RepoPath != "" {
		query.Where("repo_path IN (?)", bun.In(strings.Split(filter.RepoPath, ",")))
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
