package worker

import (
	"os"
	"path/filepath"
	"sync"

	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
)

const (
	dataDir = "/data"
	maxfilesize int64 = 4 << 30
)

type WriteResult struct {
	filename string
	bytesize int
	startAddr string
	err error
}

type WorkerJob struct {
	data []byte
	result chan *WriteResult
}


type Worker struct {
	id int
	Job chan WorkerJob
	allocator *fileallocator.FileAllocator
	currentFile string
	currentSize int
	wg *sync.WaitGroup
}

func NewWorker(id int, job chan WorkerJob, allocator *fileallocator.FileAllocator, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:id,
		Job: job,
		allocator: allocator,
		wg: wg,
	}
}

func (w *Worker) SetCurrentFile(file string) {
	w.currentFile = file
}

func (w *Worker) Run() {
	defer w.wg.Done()


	if err := w.OpenFile(); err != nil {

	}


	for job := range w.Job {

		result := w.HandleWrite(job.data)

		job.result <- result

		close(job.result)
	}


}

func (w *Worker) OpenFile() error {

	filename := w.allocator.Next()

	filePath := filepath.Join(
		dataDir,
		filename,
	);

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND, 0644)

	
	if err != nil {
		return err
	}

	defer file.Close()
	
	

	w.currentFile = filename
	w.currentSize = 0

	return nil
}

func (w *Worker) HandleWrite(data []byte) *WriteResult {

	dataSize := len(data)

	// single request datasize exceed the max file size rules
	if dataSize > int(maxfilesize) {

	}

	if w.currentSize + dataSize > int(maxfilesize) {

	}

	file, err := os.OpenFile(w.currentFile, os.O_APPEND|os.O_CREATE, 0644);

	if err != nil {
		err = file.Close()
	}

	defer file.Close()

	startingAddr := data[min(10, len(data))]

	n, err := file.Write(data)

	if err != nil {
		return &WriteResult{
			err: err,
			startAddr: string(startingAddr),
		}
	}

	
	if n != len(data) {
		return &WriteResult{
			err: err,
			startAddr: string(startingAddr),
		}
	}


	return &WriteResult{
		filename: w.currentFile,
		startAddr: string(startingAddr),
		bytesize: dataSize,
		err: nil,
	}
}