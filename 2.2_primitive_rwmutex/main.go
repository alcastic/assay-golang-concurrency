package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

type Value struct {
	v int
	// sync.RWMutex give more control over memory,
	// locks access for reading are granted unless the locks is being held for a writing
	rwmu sync.RWMutex
}

func main() {

	generateWriters := func(wg *sync.WaitGroup, value *int, locker sync.Locker) {
		defer wg.Done()
		const nroWriters = 5
		wg.Add(nroWriters)
		for i := 0; i < nroWriters; i++ {
			go func(v int) {
				defer wg.Done()
				locker.Lock()
				defer locker.Unlock()
				*value = v // Do something interesting with the shared variable
			}(i)
		}
	}

	generateReaders := func(wg *sync.WaitGroup, nroReaders int, value *int, locker sync.Locker) {
		wg.Done()
		for i := 0; i < nroReaders; i++ {
			go func() {
				locker.Lock()
				defer locker.Unlock()
				_ = *value // Do something interesting with the shared variable
			}()
		}
	}

	test := func(nroReaders int, value *Value, lockerForWrite sync.Locker, lockerForRead sync.Locker) time.Duration {
		startTest := time.Now()
		var wg sync.WaitGroup
		wg.Add(1)
		go generateWriters(&wg, &value.v, lockerForWrite)
		wg.Add(1)
		go generateReaders(&wg, nroReaders, &value.v, lockerForRead)
		wg.Wait()
		return time.Since(startTest)
	}

	v := new(Value)
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	fmt.Fprintf(tw, "Nro-Readers\tMutex\tRWMutex\n")
	for i := 0; i < 20; i++ {
		nroReaders := int(math.Pow(2, float64(i)))
		fmt.Fprintf(tw,
			"%d\t%v\t%v\n",
			nroReaders,
			test(nroReaders, v, &v.rwmu, &v.rwmu),
			test(nroReaders, v, &v.rwmu, v.rwmu.RLocker()))
	}
	tw.Flush()

}
