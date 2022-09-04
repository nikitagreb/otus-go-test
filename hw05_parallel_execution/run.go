package hw05parallelexecution

import (
	"context"
	"errors"
	"sync"
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
	taskCh := make(chan Task, len(tasks))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				select {
				case <-ctx.Done():
					return
				case task, ok := <-taskCh:
					if !ok {
						cancel()
						return
					}

					if task() != nil {
						mu.Lock()
						maxErrors--
						if maxErrors <= 0 {
							mu.Unlock()
							cancel()
							return
						}
						mu.Unlock()
					}
				default:
				}
			}
		}()
	}
	for _, task := range tasks {
		taskCh <- task
	}

	close(taskCh)
	wg.Wait()

	if maxErrors <= 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
