package dbx

import "database/sql"


var conn *sql.DB

func init() {
	var err error

	conn, err = sql.Open("sqlite", "./app.db")

	if err != nil {

	}
}

func Connection() *sql.DB {
	return conn
}

type Object struct {
	Filename string
	ObjectSize int
	StartAddress int64
	Object_key string
	Parts int
}

func InsertObject(filename string, objectSize int, startAddress int64, object_key string, parts int) error {
	query := `INSERT INTO object_mapping (file_name,object_size,start_address,object_key,parts)
	        VALUES (?, ?, ?, ?, ?)`;
	
	result, err := conn.Exec(query, filename, objectSize, startAddress, object_key, parts)

	if err != nil {

	}
	_, err = result.RowsAffected()

	if err != nil {

	}

	return nil
}

func Select(object_key string) (*Object, error) {

	var object Object;

	query := `SELECT file_name, object_size, start_address, object_key, parts WHERE object_mapping object_key = ?`;

	err := conn.QueryRow(query, object_key).Scan(
		&object.Filename,
		&object.ObjectSize,
		&object.StartAddress,
		&object.Object_key,
		&object.Parts,
	);

	if err != nil {
		return nil, err
	}

	return &object, nil
} 