package engine

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// SqlDB is the active SQLite database handle
var SqlDB *sql.DB

// ConnectSQLite opens a SQLite database file and sets the package-level handle.
func ConnectSQLite(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("sqlite open: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("sqlite ping: %w", err)
	}
	db.Exec("PRAGMA journal_mode=WAL")
	SqlDB = db
	return nil
}
