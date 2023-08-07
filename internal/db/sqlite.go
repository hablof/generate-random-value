package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func NewDB() (*sql.DB, error) {

	// f, err := os.OpenFile(filepath.Join("data", "data.db"), os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	return nil, err
	// }
	// f.Close()

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=rwc" /*filepath.Join("data",*/, "data.db")) //)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS vals(
		id INTEGER PRIMARY KEY UNIQUE,
		val VARCHAR(256) NOT NULL,
		request_id VARCHAR(256) UNIQUE);
	CREATE INDEX IF NOT EXISTS idx_req_id ON vals(request_id);`)

	if err != nil {
		return nil, err
	}

	return db, nil
}
