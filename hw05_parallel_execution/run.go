package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type ExecutionResult struct {
	mu          sync.Mutex
	maxErrors   int
	errorsCount int
}

func (r *ExecutionResult) inc() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.errorsCount++
}

func (r *ExecutionResult) isMaxErrorsReached() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.maxErrors <= 0 {
		return false
	}

	return r.errorsCount >= r.maxErrors
}

func fillTasksChannel(tasks []Task) chan Task {
	tasksChannel := make(chan Task, len(tasks))
	for i := 0; i < len(tasks); i++ {
		tasksChannel <- tasks[i]
	}
	close(tasksChannel)
	return tasksChannel
}

// Run starts tasks in goroutinesCount goroutines and stops its work when receiving maxErrors errors from tasks.
func Run(tasks []Task, goroutinesCount, maxErrors int) error {
	result := ExecutionResult{maxErrors: maxErrors}
	wg := sync.WaitGroup{}

	tasksChannel := fillTasksChannel(tasks)

	for i := 0; i < goroutinesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range tasksChannel {
				if result.isMaxErrorsReached() {
					break
				}

				err := task()
				if err != nil {
					result.inc()
				}
			}
		}()
	}

	wg.Wait()

	if result.isMaxErrorsReached() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
