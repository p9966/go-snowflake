package main

import (
	"fmt"
	"sync"

	gosnowflake "github.com/p9966/go-snowflake"
)

var wg sync.WaitGroup

// 并发生成1000个，测试是否存在重复id
func main() {
	w := gosnowflake.NewWorker(1, 1)
	count := 1000
	ch := make(chan uint64, count)
	wg.Add(count)
	defer close(ch)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			id, _ := w.NextID()
			ch <- id
		}()
	}
	wg.Wait()
	m := make(map[uint64]int)
	for i := 0; i < count; i++ {
		id := <-ch
		if _, ok := m[id]; ok {
			fmt.Printf("repeat id %d\n", id)
			return
		}
		m[id] = i
	}

	fmt.Println("All", len(m), "snowflake ID Get successed!")
}
