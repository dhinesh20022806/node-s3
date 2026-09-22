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

// 512 mb files limited
// same time 8 cores write concurrent should be done
// make sure pick the read file should be unique among other cores
// read can pick multiple thread can access it no restict on created goroutine for read operation

// return offset, size of the bytes and file name 
// how to choose the file, check file is meet the 512 mb , for now if current data where going split means
// move to new file to read those data, which means currentSize + fileSize <= 512 mb, look for 
// other remaining 7 read-write files, if not found, then  not new file
// create table for track the state of files in sql.db  like file name, file_path, state(idle, writing),
// size , permission (rw,r), timestamp
// data/a , b ...z after a1, b1...z1 , an ... zn.
// make sure that new file write by this goroutine only and only up to 8 goroutine for write.

// Limit concurrency to 8 using a buffered channel semaphore.
// var sem = make(chan struct{}, 8)

// gRPC 

func WriteObject(data []byte, filePath string) int64 {
	
	// return offset, size and filename
	
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY ,0644)
  
	 defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	offset, err := file.Seek(0, io.SeekEnd) // what is meaning of 0 here

	if err != nil {
		log.Fatal(err)
	}

	file.Write(data)


	return offset
}

func object() {

	db, err := sql.Open("sqlite", "../app.db")

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
		log.Fatalf("Failed to create table: %v", err)
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