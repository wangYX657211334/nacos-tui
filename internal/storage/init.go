package storage

import (
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"os"
	"path/filepath"
)

func NewConnection() (db *sql.DB, err error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(userHome, ".nacos-tui", "nacos-tui.db")
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// 创建表
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS audit (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        session_id varchar(8),
        base_url varchar(256),
        username varchar(256),
        password varchar(256),
        namespace varchar(256),
        request_dump TEXT,
        response_dump TEXT,
        error_message TEXT,
        time INTEGER
    );
    CREATE TABLE IF NOT EXISTS nacos_context (
        name varchar(32) PRIMARY KEY,
        url varchar(256),
        username varchar(256),
        password varchar(256),
        namespace varchar(256),
        namespace_name varchar(256)
    );
    CREATE TABLE IF NOT EXISTS system_config (
        key varchar(32) PRIMARY KEY,
        value varchar(256)
    );
    `
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
