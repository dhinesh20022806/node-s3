package workerpool

import (
	"sync"

	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
)


type WorkerPool struct {
	Workers []*worker.Worker
	queueSize int

	wg sync.WaitGroup
}

func NewWorkerPool(queueSize int, numWorker int, allocator *fileallocator.FileAllocator) *WorkerPool {
	
	pool := &WorkerPool{
		Workers: make([]*worker.Worker, numWorker),
		queueSize: queueSize,
	}

	for i := 0; i < numWorker; i++ {
		work := worker.NewWorker(i, make(chan worker.WorkerJob, queueSize), allocator, &pool.wg)
		pool.Workers = append(pool.Workers, work)
	}
	return pool
}

func (p *WorkerPool) Start() {

	for _, worker := range p.Workers {

		p.wg.Add(1)

		go worker.Run()
	}

}

func (p *WorkerPool) Shutdown() {

	for _, worker := range p.Workers {
		close(worker.Job)
	}

	p.wg.Wait()

}
