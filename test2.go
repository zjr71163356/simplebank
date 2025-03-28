package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 模拟数据源
func dataSource(out chan<- int, done <-chan struct{}) {
	// 在函数退出时关闭输出通道
	defer close(out)

	for i := 1; ; i++ {
		// 随机生成数据
		data := rand.Intn(100)

		// 使用select尝试发送数据或检查是否应该退出
		select {
		case out <- data:
			// 数据发送成功
			fmt.Printf("发送数据: %d\n", data)
		case <-done:
			// 收到退出信号
			fmt.Println("数据源收到退出信号")
			return
		}

		// 模拟处理时间
		time.Sleep(time.Millisecond * 500)
	}
}

// 模拟数据处理器
func dataProcessor(in <-chan int, processed chan<- int, done <-chan struct{}) {
	defer close(processed)

	for {
		select {
		case data, ok := <-in:
			if !ok {
				// 输入通道已关闭
				return
			}
			// 处理数据 (示例：乘以2)
			result := data * 2

			// 尝试将处理后的数据发送出去(带超时)
			select {
			case processed <- result:
				fmt.Printf("处理数据: %d -> %d\n", data, result)
			case <-time.After(300 * time.Millisecond):
				fmt.Printf("处理数据超时: %d\n", data)
			}

		case <-done:
			// 收到退出信号
			fmt.Println("处理器收到退出信号")
			return
		}
	}
}

func main() {
	// 随机数种子
	rand.Seed(time.Now().UnixNano())

	// 创建通道
	dataChan := make(chan int, 5)      // 数据通道
	processedChan := make(chan int, 5) // 处理后数据通道
	done := make(chan struct{})        // 退出信号通道

	// 启动数据源和处理器
	go dataSource(dataChan, done)
	go dataProcessor(dataChan, processedChan, done)

	// 主循环
	count := 0
	for {
		select {
		case result, ok := <-processedChan:
			if !ok {
				fmt.Println("处理通道已关闭")
				return
			}
			fmt.Printf("收到处理结果: %d\n", result)
			count++

			// 演示非阻塞读取
			select {
			case extraResult := <-processedChan:
				fmt.Printf("额外收到一个结果: %d\n", extraResult)
				count++
			default:
				// 没有更多结果可读
			}

		case <-time.After(5 * time.Second):
			// 运行5秒后退出
			fmt.Println("程序运行超时，准备退出")
			close(done) // 发送退出信号

		}

		// 处理10个数据后退出
		if count >= 10 {
			fmt.Println("已处理足够数据，准备退出")
			close(done)

			// 确保所有goroutine有时间清理
			time.Sleep(time.Second)
			return
		}
	}
}
