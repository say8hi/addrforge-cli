package worker

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type JobFunc func(workerID int) (bool, error)

func Run(workers int, job JobFunc) error {
	stop := make(chan struct{})
	wg := sync.WaitGroup{}
	var found atomic.Bool
	var counter atomic.Int64

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				tps := counter.Swap(0)
				fmt.Printf("\rTPS: %d", tps)
			}
		}
	}()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
					ok, _ := job(id)
					counter.Add(1)

					if ok && !found.Load() {
						if found.CompareAndSwap(false, true) {
							close(stop)
							return
						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
	return nil
}
