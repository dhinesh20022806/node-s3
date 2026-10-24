package dbx

import (
	"database/sql"
	"fmt"

	//    import sqlite
	_ "github.com/mattn/go-sqlite3"
)

var conn *sql.DB

func init() {
	var err error

	conn, err = sql.Open("sqlite3", "./app.db")

	if err != nil {

		panic(err)
	}
}

func Connection() *sql.DB {
	return conn
}

type Object struct {
	Filename     string
	ObjectSize   int
	StartAddress int64
	Object_key   string
	Parts        int
}

func InsertObject(filename string, objectSize int, startAddress int64, object_key string, parts int) error {

	createTableQuery := `
        CREATE TABLE IF NOT EXISTS object_mapping (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            file_name TEXT NOT NULL,
            object_size INTEGER NOT NULL,
            start_address INTEGER NOT NULL,
            object_key TEXT NOT NULL,
            parts INTEGER NOT NULL
        );
    `

	_, err := conn.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create object_mapping table: %w", err)
	}

	query := `INSERT INTO object_mapping (file_name,object_size,start_address,object_key,parts)
	        VALUES (?, ?, ?, ?, ?)`

	result, err := conn.Exec(query, filename, objectSize, startAddress, object_key, parts)

	if err != nil {

		return err
	}
	_, err = result.RowsAffected()

	if err != nil {

		return err
	}

	return nil
}

func Select(object_key string) (*Object, error) {

	var object Object

	query := `SELECT file_name, object_size, start_address, object_key, parts FROM object_mapping WHERE object_key = ?`

	err := conn.QueryRow(query, object_key).Scan(
		&object.Filename,
		&object.ObjectSize,
		&object.StartAddress,
		&object.Object_key,
		&object.Parts,
	)

	if err != nil {
		return nil, err
	}

	return &object, nil
}
