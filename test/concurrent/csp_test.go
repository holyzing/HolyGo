package concurrent

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCspBase(t *testing.T) {
	// 为了更好地编写并发程序，从设计之初Go语言就注重如何在编程语言层级上设计一个简洁安全高效的抽象模型，
	// 让开发人员专注于分解问题和组合方案，而且不用被线程管理和信号互斥这些烦琐的操作分散精力。

	// 在并发编程中，对共享资源的正确访问需要精确地控制，在目前的绝大多数语言中，都是通过加锁等线程
	// 同步方案来解决这一困难问题，而Go语言却另辟蹊径，它将共享的值通过通道传递（实际上多个独立执行的
	// 线程很少主动共享资源）。

	// 由于 mu.Lock() 和 mu.Unlock() 并不在同一个 Goroutine 中，所以也就不满足顺序一致性内存模型。
	// 同时它们也没有其他的同步事件可以参考，也就是说这两件事是可以并发的。
	// 因为可能是并发的事件，所以 main() 函数中的 mu.Unlock() 很有可能先发生，
	// 而这个时刻 mu 互斥对象还处于未加锁的状态，因而会导致运行时异常。

	/*
		var mutx sync.Mutex
		go func() {
			fmt.Println("calm down")
			mutx.Lock()
		}()
		mutx.Unlock()
		// fatal error: sync: unlock of unlocked mutex
	*/

	var mu sync.Mutex
	mu.Lock()
	go func() {
		fmt.Println("calm down")
		time.Sleep(2 * time.Second)
		mu.Unlock()
		time.Sleep(3 * time.Second)
		println("Released ！")
	}()

	mu.Lock()
	println("Blocked !")
}
