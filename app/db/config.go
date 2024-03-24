package db

import "database/sql"

type Database struct {
	Conn *sql.DB
	Cfg  DatabaseCfg
}

type DatabaseCfg struct {
	DbPath string
}
