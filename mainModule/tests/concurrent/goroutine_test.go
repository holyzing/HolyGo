package concurrent

import (
	"fmt"
	"testing"
)

/**
/**
并发编程含义比较广泛，包含多线程编程、多进程编程及分布式程序等
Go 语言通过编译器运行时（runtime），从语言上支持了并发的特性。
Go 语言的并发通过 goroutine 特性完成。
goroutine 类似于线程，但是可以根据需要创建多个 goroutine 并发工作。
goroutine 是由 Go 语言的运行时调度完成，而线程是由操作系统调度完成。
Go 语言还提供 channel 在多个 goroutine 间进行通信。
goroutine 和 channel 是 Go 语言秉承的 CSP（Communicating Sequential Process）并发模式的重要实现基础
*/

/**
TODO
make
new
type()
struct{}
*/

func TestHello(t *testing.T) {
	// 阻塞式的无缓冲的通道
	stringChan := make(chan string)
	producer := func() {
		fmt.Println("send befor ...")
		stringChan <- "message1"
		fmt.Println("received finished")
	}
	go producer()
	// consumer := <-stringChan
	// fmt.Println(consumer)

}
