package fileallocator

import (
	"sync"

	filenamegenerator "github.com/dhinesh20022806/node-s3/app/pkg/fileNameGenerator"
)

type FileAllocator struct {
	mu        sync.Mutex
	generator *filenamegenerator.FileNameGenerator
}

func NewFileAllocator(startIndex uint64) *FileAllocator {
	return &FileAllocator{
		generator: filenamegenerator.NewFileNameGenerator(startIndex),
	}
}

func (a *FileAllocator) Next() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.generator.Next()
}
