package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	id       int
	workChan chan string
	killChan chan struct{}
}

type WorkerPool struct {
	workers  map[int]*Worker
	poolChan chan string
	shutdown chan struct{}
	count    int
	wg       sync.WaitGroup
	mutex    sync.Mutex
}

func NewWorkerPool() *WorkerPool {
	wp := &WorkerPool{
		workers:  make(map[int]*Worker),
		poolChan: make(chan string),
		shutdown: make(chan struct{}),
		count:    0,
		wg:       sync.WaitGroup{},
		mutex:    sync.Mutex{},
	}

	return wp
}

func (wp *WorkerPool) AddWorker() {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	worker := &Worker{
		id:       wp.count + 1,
		workChan: wp.poolChan,
		killChan: make(chan struct{}),
	}

	wp.wg.Add(1)
	go worker.runWorker(&wp.wg)

	wp.workers[worker.id] = worker
	wp.count++
}

func (wp *WorkerPool) RemoveWorker(id int) {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if worker, ok := wp.workers[id]; ok {
		close(worker.killChan)
		delete(wp.workers, id)
		return
	}
}

func (wp *WorkerPool) Shutdown() {
	wp.mutex.Lock()
	close(wp.shutdown)

	for _, worker := range wp.workers {
		close(worker.killChan)
	}
	close(wp.poolChan)
	wp.mutex.Unlock()

	wp.wg.Wait()
}

func (wp *WorkerPool) SendTask(task string) error {
	select {
	case wp.poolChan <- fmt.Sprintf("Task %s", task):
		return nil
	case <-wp.shutdown:
		return errors.New("worker-pool shutdowned")
	}
}

func (worker *Worker) runWorker(wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range worker.workChan {
		done := make(chan struct{})
		go func() {
			dataProcessing(worker.id, task)
			close(done)
		}()

		select {
		case <-done:
			fmt.Printf("Worker %d done\n", worker.id)
		case <-worker.killChan:
			fmt.Printf("Worker %d killed during task\n", worker.id)
			return
		}
	}
}

func dataProcessing(wId int, data string) {
	fmt.Printf("Worker %d started: %s\n", wId, data)
	time.Sleep(1 * time.Second)
}

func main() {
	wp := NewWorkerPool()
	wp.AddWorker()
	wp.AddWorker()
	wp.AddWorker()

	go func() {
		for i := 0; i < 10; i++ {
			err := wp.SendTask(strconv.Itoa(i))
			if err != nil {
				return
			}
		}
	}()
	time.Sleep(time.Millisecond * 500)
	wp.RemoveWorker(1)
	wp.AddWorker()
	time.Sleep(time.Millisecond * 500)
	wp.Shutdown()
}
