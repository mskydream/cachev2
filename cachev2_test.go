package cachev2

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestRace(t *testing.T) {
	c := New()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i % 5)

			for j := 0; j < 1000; j++ {
				c.Set(key, j, time.Millisecond)
				c.Get(key)
				c.Delete(key)
			}
		}(i)
	}

	wg.Wait()
}
