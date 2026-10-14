package worker

import "sync"


type Worker struct {
	id int
	job chan string
	allocator string
	currentFile string
	currentSize int
	wg *sync.WaitGroup
}