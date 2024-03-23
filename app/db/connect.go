package db

import (
	"database/sql"
	"fmt"
)

func (d *Database) ConnectDatabase() error {
	db, err := sql.Open("sqlite3", d.Cfg.DbPath)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	d.Conn = db
	return nil
}
