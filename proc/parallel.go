// Package proc provides parallel processing utilities.
package proc

import (
	"runtime"
	"sync"
)

func Parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}

	procs := min(runtime.GOMAXPROCS(0), count)

	c := make(chan int, count)
	for i := range stop {
		c <- start + i
	}
	close(c)

	var wg sync.WaitGroup
	for range procs {
		wg.Go(func() { fn(c) })
	}
	wg.Wait()
}
