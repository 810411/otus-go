package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in N goroutines and stops its work when receiving M errors from tasks
func Run(tasks []Task, n int, m int) error {
	errLimit := m
	taskCh := make(chan Task, len(tasks))
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			worker(taskCh, &mu, &errLimit)
		}()
	}

	for _, t := range tasks {
		taskCh <- t
	}

	close(taskCh)
	wg.Wait()

	if errLimit < 0 {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func worker(taskCh <-chan Task, mu sync.Locker, errLimit *int) {
	stop := false
	for t := range taskCh {
		err := t()
		mu.Lock()
		if *errLimit < 1 {
			stop = true
		}
		if err != nil {
			*errLimit--
		}
		mu.Unlock()
		if stop {
			return
		}
	}
}
