# go-snowflake
基于雪花算法的`Golang`分布式生成器

## 特性
❤ 零配置，开箱即用  
🚀 高并发分布式系统环境下ID不重复  
🧭 生成效率高  
🐱‍🐉 不依赖于第三方的库或者中间件  
✔ 生成的id具备时序性和唯一性  

## 安装
```go
go get github.com/p9966/go-snowflake
```

## 用法
```go
w := gosnowflake.NewWorker(1, 1)
id, _ := w.NextID()
```

```go
var wg sync.WaitGroup

// 并发生成1000个，测试是否存在重复id
func main() {
	w := gosnowflake.NewWorker(1, 1)
	id, _ := w.NextID()
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
```
