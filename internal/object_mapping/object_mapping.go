package objectmapping

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"uuid"
)


type ObjectMapping struct {
	ID string
	ObjectId string
	FileName string
	StartOffset byte
	ObjectSize int
}


func object() {

	db, err := sql.Open("sqlite", "./app.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()


	createTableSQL := `CREATE TABLE IF NOT EXISTS object_mapping (
	               id TEXT PRIMARY KEY,
				   object_id TEXT,
				   file_name TEXT,
				   start_offset TEXT,
				   object_size INT
	);`

	_, err = db.Exec(createTableSQL)

	if err != nil {
		log.Fatal("Failed to create table: %v", err)
	}

	fmt.Println("Table create or verified successfully.")


	blob, err := os.ReadFile("images.webp")

	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}


	size := len(blob)

	beginingAddress := blob[:10]

	fmt.Printf("size of the blob  %d and %s \n", size, beginingAddress)

	file, err := os.OpenFile("outputs", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}

	defer file.Close()

	if _, err = file.Write(blob); err != nil {
		log.Fatalf("Failed to write Blob to file: %v", err)
		
	}

	newUUID := uuid.New().String()

	insertSQL := `INSERT INTO object_mapping (id, object_id, file_name, start_offset, object_size)
	            values (?, ?, ?, ?, ?);`
	
	if _, err = db.Exec(insertSQL, newUUID, newUUID, "outputs", beginingAddress, size); err != nil {
		log.Fatal(err)
	}
	// err = os.WriteFile("outputs", blob, 0644)

// /

    offset, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(offset)


	bringBackEle()
	


}



func bringBackEle() {

	size := 42286

	file, err := os.Open("output")

	if err != nil {
		log.Fatalf("Failed to open file: %v", err)

	}

	defer file.Close()

	// file.ReadAt()

	position, err := file.Seek(0, io.SeekStart)



	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("MOved to absolute bytes posisiton: %d\n", position)

	buffer := make([]byte, size)

	_, err = file.Read(buffer)

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("output.webp", buffer, 0644)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("completed recoved image from bytes")



}