package objectmapping

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"uuid"

	_ "modernc.org/sqlite"
)


type ObjectMapping struct {
	ID string
	ObjectId string
	FileName string
	StartOffset int64
	ObjectSize int64
}

type PointerMeta struct {
	pointer int
	alphabetPointer int
}

type Permission string

const (
	ReadAndWrite Permission = "rw"
	ReadOnly Permission = "r"
)

type InputData struct {
	key string
	data []byte
}

type Queue struct {
	queue []InputData
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

// how to handle internal fragmentation

// who is reponsible for creating file


func CreateFile() string {
	 alphabet := "abcdefghijklmnopqrstuvwxyz"

	 data, err := readPointer()


	 if err != nil {
		// fmt.Printf("data: %d %d\n", data.pointer, data.alphabetPointer)

		if errors.Is(err, sql.ErrNoRows) {
			data = &PointerMeta{pointer: 0, alphabetPointer: 0}
		} else {
			log.Fatalf("err %v", err)
		}
	 }

	pointer := data.pointer 
	alphabet_pointer := data.alphabetPointer;

	 fileName := fmt.Sprintf("../../data/%s%d", string(alphabet[alphabet_pointer]), pointer)

	 os.Create(fileName)
	 alphabet_pointer += 1
	 if alphabet_pointer == 26 {
		pointer += 1
		alphabet_pointer = 0
	 }
	 updatePointers(pointer, alphabet_pointer)
	 return fileName;
}

func main(){


	queue := []string{"dog", "cat", "bat", "hat"}

	go func() {
		for {
			for i := 0; i < len(queue); i++ {

			go func(data string) {

               fmt.Print(data)

			}(queue[i])

			}
			
		}
	}()

	go func() {

		queue = append(queue, "")


	}()

	//  cors := runtime.NumCPU()

	//  var sem = make(chan struct{}, cors)


	 // file name should be data/a data/b ... data/z 
	 // data/a1 data/b1 ... data/z1




	 createTableSQL := `CREATE TABLE IF NOT EXISTS object_mapping (
	               id TEXT PRIMARY KEY,
				   object_id TEXT,
				   file_id TEXT,
				   start_offset TEXT,
				   object_size INTEGER,
				   deleted_at TIMESTAMP,
				   FOREIGN KEY (file_id) REFERENCES file_detail(id) ON DELETE CACADE
	);`;

	 createSQLFile := `CREATE TABLE IF NOT EXISTS file_detail (
	               id TEXT PRIMARY KEY,
				   file_path TEXT,
				   state TEXT CHECK (state IN ('idle', 'write')),
				   size INTEGER,
				   permission TEXT CHECK (permission IN ('r', 'rw')), 
				   timestamp DATATIME DEFAULT CURRENT_TIMESTAMP`;
	 
	createPointer := `CREATE TABLE IF NOT EXISTS pointer_meta (
	id TEXT PRIMARY KEY,
	pointer INTEGER DEFAULT 0 NOT NULL,
	alphabet_pointer INTEGER DEFAULT 0 NOT NULL,
	timestamp DATATIME DEFAULT CURRENT_TIMESTAMP
	)`;

	_, err := db.Exec(createTableSQL);
	if err != nil {
		log.Fatalf("Failed to create table: %v:", err)
	}
	_, err = db.Exec(createSQLFile)
	if err != nil {
		log.Fatalf("Failed to Create Sql File Table: %v", err)
	}
	
	_, err = db.Exec(createPointer)

	if err != nil {
		log.Fatalf("Failed to Create pointer Table: %v", err)
	}

	fmt.Printf("Tables are Created!\n")

	// when to change the file from rw to r ? when file is inability to store data ? it's works but how to handle 
	// that internal fragramentation ?

	// here the idea , instead of checking 512 mb large for each file , when about to exceed the file size 
	// make the last data append to that file and mark it's readonly, instead of capcity now it's miniumum 
	// size requirement to create new file , is this create external fragementation ?



	// when to increase the pointer , when there is need to create new file after wards

	// go routine for 8 max 

	// idle + more space

	// elezabeth olsen, scarlett johansson, ... 6

	// check if 8 eight file are writing , if no create file move to writeObject


	count := readWriteFileCount()

	// plus in demand 
	// minium eight read-and-write file exists 

	if count < 8 {
		// TODO: create new file
		// pass it

		filePath := CreateFile()

		WriteObject([]byte("Hello, World"), filePath);
	}

	
}

// pick any available non active file // write a function

func activeFile() int {
	query := `SELECT COUNT(*) FROM file_detail WHERE state = 'write'`;

	var count int;
	if err := db.QueryRow(query).Scan(&count); err != nil {
		log.Fatalln(err)
	}
	return count
}

func readWriteFileCount() int {
	query := `SELECT COUNT(*) FROM file_detail WHERE permission = "rw";`

	var count int;

	if err := db.QueryRow(query).Scan(&count); err != nil {
		log.Fatalln(err)
	}
	return count
}


type AiRequest struct {
	Messages string `json:message`
}

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite", "../../app.db")
	if err != nil {
		return err
	}

	return db.Ping()
}

var MAX_FILE_SIZE = 512 * 1024 * 1024

func isAboveRequiredFileSize(size int) bool {
	return size <= MAX_FILE_SIZE
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}


func ReadObject(objectId string) []byte {

	var obj ObjectMapping;

	fmt.Print(runtime.NumCPU())

readSQL := `
	SELECT id, object_id, file_name, start_offset, object_size
	FROM object_mapping
	WHERE object_id = ?
`

err := db.QueryRow(readSQL, objectId).Scan(
	&obj.ID,
	&obj.ObjectId,
	&obj.FileName,
	&obj.StartOffset,
	&obj.ObjectSize,
)

	if err != nil {
		log.Fatalf("error at scan %v", err)
	}

	file, err := os.Open(obj.FileName)

	if err != nil {
		log.Fatalf("error at scan %v", err)
	}

	_, err = file.Seek(int64(obj.StartOffset), io.SeekStart)

	buffer := make([]byte, obj.ObjectSize)

	bufferRead, err := file.Read(buffer)

	if err != nil {
		log.Fatalf("buffer read error %v", err)
	}

	return buffer[:bufferRead]
}


func WriteRepo(offset int64, fileName string, objectKey string, size int) string {

	newUUID := uuid.New().String()

	insertSQL := `INSERT INTO object_mapping (id, object_id, file_name, start_offset, object_size)
	            values (?, ?, ?, ?, ?);`
	
	if _, err := db.Exec(insertSQL, newUUID, newUUID, fileName, offset, size); err != nil {
		log.Fatal(err)
	}

	return newUUID
}
func updatePointers(pointer int, alphabet_pointer int) {
	update_query := `UPDATE pointer_meta SET pointer = ?, alphabet_pointer = ?`
	if _, err := db.Exec(update_query, pointer, alphabet_pointer); err != nil {
		log.Fatal(err)
	}
}

func InsertPointer() error {

	// var obj PointerMeta 

	exec_query := `INSERT INTO pointer_meta (pointer, alphabet_pointer) VALUES
	                (?, ?)`;
	_, err := db.Exec(exec_query, 0, 0)


	if err != nil {

		return err
	}
	return nil

}

func readPointer() (*PointerMeta, error) {

	var obj PointerMeta

	read_query := `SELECT pointer, alphabet_pointer FROM pointer_meta`

	err := db.QueryRow(read_query).Scan(
		&obj.pointer,
		&obj.alphabetPointer,
	)

	if err != nil {
		return nil, err
	}
	return &obj, nil
}



func InsertFileRepo(filename string, permission Permission) string {

	newUUID := uuid.New().String()

	insertSQL := `INSERT INTO file_detail (id, file_path, state, size, permission)
	              value (?, ?, ?, ?, ?)`;
	
    if _, err := db.Exec(insertSQL, newUUID, filename, "idle", 0, permission ); err != nil {
		log.Fatalf("Error occured when insert new file %v", err)
	}

	return newUUID
}

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