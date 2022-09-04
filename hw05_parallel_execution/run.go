package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkers      = errors.New("the number of workers must be greater than zero")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// maxErrors <= 0 then there is no limit on errors
func Run(tasks []Task, workers int, maxErrors int) error {
	if workers <= 0 {
		return ErrInvalidWorkers
	}
	if maxErrors <= 0 {
		maxErrors = len(tasks) + 1
	}
	taskChannel := make(chan Task)
	var errorsCount int32
	wg := sync.WaitGroup{}
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for task := range taskChannel {
				err := task()
				if err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errorsCount) >= int32(maxErrors) {
			break
		}
		taskChannel <- task
	}

	close(taskChannel)
	wg.Wait()

	if errorsCount >= int32(maxErrors) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
