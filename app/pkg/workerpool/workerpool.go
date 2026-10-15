package workerpool

import (
	"sync"

	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
)


type WorkerPool struct {
	workers []*worker.Worker
	queueSize int

	wg sync.WaitGroup
}

func NewWorkerPool(queueSize int, numWorker int, allocator *fileallocator.FileAllocator) *WorkerPool {
	
	pool := &WorkerPool{
		workers: make([]*worker.Worker, numWorker),
		queueSize: queueSize,
	}

	for i := 0; i < numWorker; i++ {
		work := worker.NewWorker(i, make(chan worker.WorkerJob, queueSize), allocator, &pool.wg)
		pool.workers = append(pool.workers, work)
	}
	return pool
}

func (p *WorkerPool) Start() {

	for _, worker := range p.workers {

		p.wg.Add(1)

		go worker.Run()
	}

}
