// Package tests Benchmark test in
package tests

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/dhinesh20022806/node-s3/app/pkg/dispatcher"
	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
	"github.com/dhinesh20022806/node-s3/app/pkg/workerpool"
)

// Single Threaded Write Object

var mu sync.Mutex

func setupDataDir(b *testing.B) string {
	b.Helper()

	dir, err := filepath.Abs("./data")
	if err != nil {
		b.Fatalf("failed to resolve data directory: %v", err)
	}

	if err := os.RemoveAll(dir); err != nil {
		b.Fatalf("failed to remove data directory: %v", err)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		b.Fatalf("failed to create data directory: %v", err)
	}

	b.Cleanup(func() {
		if err := os.RemoveAll(dir); err != nil {
			b.Errorf("failed to remove data directory: %v", err)
		}
	})

	return dir
}

func SingleThreadWriteObject(data []byte) {

	mu.Lock()
	defer mu.Unlock()

	file, err := os.OpenFile("./abc", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {

	}

	file.Write(data)
}

func BenchmarkSingleThread(b *testing.B) {

	resp, err := http.Get("https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png")

	if err != nil {
		b.Fatalf("failed to download image: %v", err)
	}

	defer resp.Body.Close()

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Fatalf("failed to read response: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		SingleThreadWriteObject(buffer)
	}
}

func BenchmarkMultiThread(b *testing.B) {

	dataDir := setupDataDir(b)

	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)

	resp, err := http.Get("https://go.dev/blog/go-brand/Go-Logo/PNG/Go-Logo_LightBlue.png")

	if err != nil {
		b.Fatalf("failed to download image: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.Fatalf("unexpected http requests: %s", resp.Status)
	}

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		b.Fatalf("failed to read response: %v", err)
	}

	fileAlloc := fileallocator.NewFileAllocator(0)

	no_worker := runtime.GOMAXPROCS(0)

	workerPool := workerpool.NewWorkerPool(8, no_worker, fileAlloc)
	workerPool.Start()
	dispatcher.NewDispatcher(workerPool)
	_ = dataDir
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {

		job := worker.WorkerJob{
			Data:   buffer,
			Result: make(chan worker.WriteResult, 1),
		}

		ctx := context.Background()

		dispatcher.DisPatcher.Submit(ctx, job)

		select {

		case result, ok := <-job.Result:

			if !ok {
				b.Fatalf("job result channel was closed ")
			}

			_ = result.Filename

		case <-ctx.Done():
			b.Fatalf("job context cancelled: %v", ctx.Err().Error())

		}
	}

}
