package worker

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type JobFunc func(workerID int) (bool, error)

type Result struct {
	Found    bool
	Data     interface{}
	Error    error
	WorkerID int
}

func Run(workers int, job JobFunc) error {
	return RunWithResult(workers, job, nil)
}

func RunWithResult(workers int, job JobFunc, onResult func(Result)) error {
	stop := make(chan struct{})
	resultChan := make(chan Result, workers)
	wg := sync.WaitGroup{}

	var found atomic.Bool
	var counter atomic.Int64
	var total atomic.Int64
	var errors atomic.Int64

	// Statistics goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				tps := counter.Swap(0)
				fmt.Printf("\rTPS: %-10d | Total: %-10d | Errors: %-5d",
					tps, total.Load(), errors.Load())
			}
		}
	}()

	// Result handler goroutine
	go func() {
		for result := range resultChan {
			if result.Error != nil {
				errors.Add(1)
				log.Printf("Worker %d error: %v", result.WorkerID, result.Error)
				continue
			}

			if result.Found && !found.Load() {
				if found.CompareAndSwap(false, true) {
					if onResult != nil {
						onResult(result)
					}
					close(stop)
					return
				}
			}
		}
	}()

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ok, err := job(id)
					counter.Add(1)
					total.Add(1)

					// Send result through channel
					select {
					case resultChan <- Result{
						Found:    ok,
						Error:    err,
						WorkerID: id,
					}:
					case <-stop:
						return
					}

					// Break early if found or error is critical
					if ok || (err != nil && isCriticalError(err)) {
						return
					}
				}
			}
		}(i)
	}

	wg.Wait()
	close(resultChan)
	fmt.Println() // New line after progress

	return nil
}

func isCriticalError(err error) bool {
	return false
}
