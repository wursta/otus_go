package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in goroutinesCount goroutines and stops its work when receiving maxErrors errors from tasks.
func Run(tasks []Task, goroutinesCount, maxErrors int) error {
	errorsCount := 0
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 1; i <= len(tasks); i++ {
		// Если количество ошибок уже превысило допустимое количество, ожидаем завершения оставшихся горутин и выходим.
		mu.Lock()
		if isMaxErrorsReached(errorsCount, maxErrors) {
			mu.Unlock()
			wg.Wait()
			break
		}
		mu.Unlock()

		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			err := tasks[i-1]()
			if err != nil {
				mu.Lock()
				errorsCount++
				mu.Unlock()
			}
		}(i)

		// Если мы запустили "пороговое" количество горутин, то дожидаемся их выполнения
		if i%goroutinesCount == 0 {
			wg.Wait()
		}
	}

	wg.Wait()

	mu.Lock()
	if isMaxErrorsReached(errorsCount, maxErrors) {
		mu.Unlock()
		return ErrErrorsLimitExceeded
	}
	mu.Unlock()

	return nil
}

func isMaxErrorsReached(errorsCount, maxErrors int) bool {
	if maxErrors <= 0 {
		return false
	}

	return errorsCount >= maxErrors
}
