package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// 强制使用单线程执行，使调度效果更明显
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	ch := make(chan int) // 无缓冲通道

	// 记录并打印当前goroutine ID
	printGID := func(prefix string) {
		// 获取goroutine ID的小技巧
		var buf [64]byte
		n := runtime.Stack(buf[:], false)
		id := string(buf[:n])
		fmt.Printf("%s: %s\n", prefix, id)
	}

	// 生产者goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			printGID(fmt.Sprintf("生产者 (发送前 %d)", i))
			ch <- i // 这里会阻塞并触发调度
			time.Sleep(time.Second * 2)
			printGID(fmt.Sprintf("生产者 (发送后 %d)", i))
			// 添加小延迟以便更容易观察到调度效果
			time.Sleep(time.Second * 2)
		}
		close(ch)
	}()

	// 消费者goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for val := range ch {
			time.Sleep(time.Second * 2)
			printGID(fmt.Sprintf("消费者 (接收到 %d)", val))
			// 添加小延迟以便更容易观察到调度效果

		}
	}()

	// 主goroutine
	printGID("主函数开始")
	wg.Wait()
	printGID("主函数结束")
}
