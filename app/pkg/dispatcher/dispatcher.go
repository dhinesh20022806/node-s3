package dispatcher

import (
	"context"
	"runtime"
	"sync/atomic"

	fileallocator "github.com/dhinesh20022806/node-s3/app/pkg/fileAllocator"
	"github.com/dhinesh20022806/node-s3/app/pkg/worker"
	"github.com/dhinesh20022806/node-s3/app/pkg/workerpool"
)

var DisPatcher *Dispatcher

func init() {

	fileAlloc := fileallocator.NewFileAllocator(0)

	no_worker := runtime.GOMAXPROCS(0)
	worker_p := workerpool.NewWorkerPool(8, no_worker, fileAlloc)
	DisPatcher = NewDispatcher(worker_p)
	worker_p.Start()

}

type Dispatcher struct {
	pool *workerpool.WorkerPool

	nextWorker atomic.Uint64
}

func NewDispatcher(pool *workerpool.WorkerPool) *Dispatcher {

	return &Dispatcher{
		pool: pool,
	}
}

func (d *Dispatcher) Submit(ctx context.Context, job worker.WorkerJob) error {

	workerCount := len(d.pool.Workers)

	index := d.nextWorker.Add(1) - 1

	work := d.pool.Workers[index%uint64(workerCount)]

	select {

	case work.Job <- job:
		return nil

	case <-ctx.Done():

		return ctx.Err()
	}
}
