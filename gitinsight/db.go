package gitinsight

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/extra/bundebug"
)

var gdb *bun.DB

func ResetDb(config *Config) error {
	err := ResetCommit()
	if err != nil {
		return err
	}
	return nil
}
func InitDb() error {
	err := InitCommit()
	if err != nil {
		return err
	}
	return nil
}

func OpenDb(typ string, dsn string) error {

	// Open database connection
	sqldb, err := sql.Open(typ, dsn)
	if err != nil {
		return err
	}

	var db *bun.DB
	if strings.Contains(typ, "mysql") {
		db = bun.NewDB(sqldb, mysqldialect.New())
	} else if strings.Contains(typ, "sqlite") {
		db = bun.NewDB(sqldb, sqlitedialect.New())
	} else {
		return errors.New("unsupported database type: " + typ)
	}

	// Add query debugging (optional)
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
	gdb = db

	return nil
}

func CloseDb() error {
	if gdb == nil {
		return nil
	}
	return gdb.Close()
}
