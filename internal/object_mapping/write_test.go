package objectmapping

import (
	"bytes"
	_ "database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"uuid"
)



func TestWriteObject(t *testing.T) {

	if err := InitDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// expect this method to write 

	// file, err := os.OpenFile("../../images.webp", os.O_RDONLY ,0644)

	blob, err := os.ReadFile("../../images.webp")

	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	
	writtenOffset := WriteObject(blob, "../../outputs")

	size := len(blob)

    
	file, err := os.Open("../../outputs")

	if err != nil {
		log.Fatalf("Failed to read image: %v", err)
	}

	offset, err := file.Seek(writtenOffset, io.SeekStart)

	if err != nil {
		log.Fatalf("Failed to Seek Error %v", err)
	}

	buffer := make([]byte, size)

	bufferRead, err := file.Read(buffer)

	if err != nil && err != io.EOF {
		log.Fatalf("Failed to read %v", err)
	}

	if offset != writtenOffset {
		t.Errorf("offset and writternOffset doesn't equal")
	}

	if !bytes.Equal(blob, buffer[:bufferRead]){
         t.Errorf("different content")
	}

	newUUid := uuid.New().String()

	key := WriteRepo(writtenOffset, "../../outputs", newUUid, size)


	expectedBytes := ReadObject(key)

	if !bytes.Equal(blob, expectedBytes){
         t.Errorf("different content after ReadObject")
	}

      
}

func TestCreateFile(t *testing.T) {

	if err := InitDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

		createPointer := `CREATE TABLE IF NOT EXISTS pointer_meta (
	id TEXT PRIMARY KEY,
	pointer INTEGER DEFAULT 0 NOT NULL,
	alphabet_pointer INTEGER DEFAULT 0 NOT NULL,
	timestamp DATATIME DEFAULT CURRENT_TIMESTAMP
	)`;

	db.Exec(createPointer)

	filename :=  CreateFile()

	fmt.Printf("the filename is %s\n",filename)

}