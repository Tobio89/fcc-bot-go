package db

import (
	"fmt"
)

func (d *Database) ConfigureDatabase() error {
	pragmaConfig := `
		PRAGMA busy_timeout = 5000;
		PRAGMA foreign_keys = ON;
		PRAGMA journal_mode = WAL;
	`
	_, err := d.Conn.Exec(pragmaConfig)
	if err != nil {
		return fmt.Errorf("failed configuring database: %w", err)
	}

	return nil
}
