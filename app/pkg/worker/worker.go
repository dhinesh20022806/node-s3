package worker

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
)

const (
	dataDir           = "./data"
	maxfilesize int64 = 4 << 30
)

type WriteResult struct {
	Filename  string
	Bytesize  int
	StartAddr int64
	Err       error
}

type WorkerJob struct {
	Data   []byte
	Result chan WriteResult
}

type Worker struct {
	id          int
	Job         chan WorkerJob
	allocator   *fileallocator.FileAllocator
	currentFile string
	currentSize int
	wg          *sync.WaitGroup
}

func NewWorker(id int, job chan WorkerJob, allocator *fileallocator.FileAllocator, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:        id,
		Job:       job,
		allocator: allocator,
		wg:        wg,
	}
}

func (w *Worker) SetCurrentFile(file string) {
	w.currentFile = file
}

func (w *Worker) Run() {
	fmt.Printf("w = %v\n", w)
	defer w.wg.Done()

	if err := w.OpenFile(); err != nil {

		fmt.Printf("file Open err = %v\n", err)
	}

	log.Printf(
		"worker %d started -> %s",
		w.id,
		w.currentFile,
	)

	for job := range w.Job {

		result := w.HandleWrite(job.Data)

		slog.Info("Result returned", "result fileName", result.Filename)
		job.Result <- result
		slog.Info("Returned to channel")

		close(job.Result)
	}

	log.Printf(
		"worker %d stopped",
		w.id,
	)

}

func (w *Worker) OpenFile() error {

	filename := w.allocator.Next()

	filePath := filepath.Join(
		dataDir,
		filename,
	)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	w.currentFile = filePath
	w.currentSize = 0

	log.Printf(
		"worker %d -> new file %s",
		w.id,
		filePath,
	)

	return nil
}

func (w *Worker) HandleWrite(data []byte) WriteResult {

	dataSize := len(data)

	// single request datasize exceed the max file size rules
	if dataSize > int(maxfilesize) {

	}

	if w.currentSize+dataSize > int(maxfilesize) {

		err := w.OpenFile()
		slog.Info("openNewFile", "newFileName", w.currentFile)

		if err != nil {
			fmt.Printf("Open failed: err = %v", err)
		}

	}
	slog.Info("file descriptor", "filename", w.currentFile)

	file, err := os.OpenFile(w.currentFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {

		slog.Error("error openingfile", "error", err)

		return WriteResult{
			Err: err,
		}
	}

	defer func() {

		err = file.Close()

		fmt.Printf("file closed err = %v", err)

		if err != nil {
			fmt.Printf("err = %v \n", err)
		}

	}()

	startingAddr, err := file.Seek(0, io.SeekCurrent)

	// startingAddr := int64(0)
	// err = nil

	if err != nil {
		panic(err)
	}

	n, err := file.Write(data)

	if err != nil {
		slog.Error("error writing file", "err", err)
		return WriteResult{
			Err:       err,
			StartAddr: startingAddr,
		}
	}

	if n != len(data) {
		return WriteResult{
			Err:       err,
			StartAddr: startingAddr,
		}
	}

	slog.Info("#data written successfully!", "filename", w.currentFile, "bytesize", dataSize)

	return WriteResult{
		Filename:  w.currentFile,
		StartAddr: startingAddr,
		Bytesize:  dataSize,
		Err:       nil,
	}
}
