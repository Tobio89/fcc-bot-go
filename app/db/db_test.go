package db

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestDatabaseOperations(t *testing.T) {
	testDBPath := "../../_store/test_db.sqlite"
	database := &Database{Cfg: DatabaseCfg{DbPath: testDBPath}}

	err := database.ConnectDatabase()
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	err = database.ConfigureDatabase()
	if err != nil {
		t.Fatalf("Error configuring database: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(testDBPath)
		t.Log("Cleanup complete")
	})

	t.Run("TestConfiguration", func(t *testing.T) {
		testConfiguration(t, database.Conn)
	})
	t.Run("TestCreateTable", func(t *testing.T) {
		testCreateTable(t, database.Conn)
	})

	t.Run("TestInsert", func(t *testing.T) {
		testInsert(t, database.Conn)
	})

	t.Run("TestSelectCount", func(t *testing.T) {
		testSelectCount(t, database.Conn)
	})

	t.Run("TestSelectValue", func(t *testing.T) {
		testSelectValue(t, database.Conn)
	})

	t.Run("TestDelete", func(t *testing.T) {
		testDeleteValue(t, database.Conn)
	})

	t.Run("TestCloseConnection", func(t *testing.T) {
		testCloseConnection(t, database.Conn)
	})
}

func testConfiguration(t *testing.T, db *sql.DB) {
	var got string
	want := "5000"

	err := db.QueryRow("PRAGMA busy_timeout").Scan(&got)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if got != want {
		t.Errorf("unexpected busy timeout: got %s, want %s", got, want)
	}
}

func testCreateTable(t *testing.T, db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)")
	if err != nil {
		t.Fatalf("failed to create table: %v", err)
	}
}

func testInsert(t *testing.T, db *sql.DB) {
	_, err := db.Exec("INSERT INTO users (name) VALUES ('Test McTesterton')")
	if err != nil {
		t.Fatalf("failed to insert into table: %v", err)
	}
}

func testSelectCount(t *testing.T, db *sql.DB) {
	var got int
	want := 1

	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&got)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if got != want {
		t.Errorf("unexpected number of rows in users table: got %d, want %d", got, want)
	}
}

func testSelectValue(t *testing.T, db *sql.DB) {
	var got string
	want := "Test McTesterton"

	err := db.QueryRow("SELECT name FROM users").Scan(&got)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if got != want {
		t.Errorf("unexpected name in users table query: got %s, want %s", got, want)
	}
}

func testDeleteValue(t *testing.T, db *sql.DB) {
	_, err := db.Exec("DELETE FROM users WHERE name = 'Test McTesterton'")
	if err != nil {
		t.Fatalf("failed to delete from table: %v", err)
	}

	var got int
	want := 0
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = 'Test McTesterton'").Scan(&got)
	if err != nil {
		t.Fatalf("failed to query database: %v", err)
	}
	if got != want {
		t.Errorf("unexpected number of rows in users table: got %d, want %d", got, want)
	}
}

func testCloseConnection(t *testing.T, db *sql.DB) {
	err := db.Close()
	if err != nil {
		t.Fatalf("failed to close database connection: %v", err)
	}

	err = db.Ping()
	if err == nil {
		t.Errorf("expected database connection to be closed, but it's still active")
	}
}
