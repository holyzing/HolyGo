package concurrent

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/**
  并发编程含义比较广泛，包含多线程编程、多进程编程及分布式程序等
  Go 语言通过编译器运行时（runtime），从语言上支持了并发的特性。
  Go 语言的并发通过 goroutine 特性完成。
  goroutine 类似于线程，但是可以根据需要创建多个 goroutine 并发工作。
  goroutine 是由 Go 语言的运行时调度完成，而线程是由操作系统调度完成。
  Go 语言还提供 channel 在多个 goroutine 间进行通信。
  goroutine 和 channel 是 Go 语言秉承的 CSP（Communicating Sequential Process）并发模式的重要实现基础。

  Go语言的并发机制运用起来非常简便，在启动并发的方式上直接添加了语言级的关键字就可以实现，
  同时实现了自动垃圾回收机制，和其他编程语言相比更加轻量。

  所有 goroutine 在 main() 函数结束时会一同结束。

  goroutine 虽然类似于线程概念，但是从调度性能上没有线程细致，
  而细致程度取决于 Go 程序的 goroutine 调度器的实现和运行环境。

  终止 goroutine 的最好方法就是 goroutine 对应的函数的自然返回，即return。
  虽然可以用 golang.org/x/net/context 包进行 goroutine 生命期深度控制，
  但这种方法仍然处于内部试验阶段，并不是官方推荐的特性。
  截止 Go 1.9 版本，暂时没有标准接口获取 goroutine 的 ID。


  NOTE 当所有的协程处于阻塞的状态后, 所有协程会处于 asleep 状态,会引发运行时的 panic
*/

/**
TODO
make
new
type()
struct{}
*/

/**
// NOTE C语言 共享内存式的并发

#include <stdio.h>
#include <stdlib.h>
#include <pthread.h>
void *count();
pthread_mutex_t mutex1 = PTHREAD_MUTEX_INITIALIZER;
int counter = 0;
int main()
{
    int rc1, rc2;
    pthread_t thread1, thread2;
    // 创建线程，每个线程独立执行函数functionC
    if((rc1 = pthread_create(&thread1, NULL, &count, NULL)))
    {
        printf("Thread creation failed: %d\n", rc1);
    }
    if((rc2 = pthread_create(&thread2, NULL, &count, NULL)))
    {
        printf("Thread creation failed: %d\n", rc2);
    }
    // 等待所有线程执行完毕
    pthread_join( thread1, NULL);
    pthread_join( thread2, NULL);
    exit(0);
}
void *count()
{
    pthread_mutex_lock( &mutex1 );
    counter++;
    printf("Counter value: %d\n",counter);
    pthread_mutex_unlock( &mutex1 );
}

*/

func TestHello(t *testing.T) {
	// NOTE 共享内存式的通信
	var count int = 0
	counter := func(lock *sync.Mutex) {
		lock.Lock()
		count++
		fmt.Println(count)
		lock.Unlock()
	}

	lock := &sync.Mutex{} //  互斥；互斥元，互斥体；互斥量
	for i := 0; i < 10; i++ {
		go counter(lock)
	}

	for {
		lock.Lock()
		c := count
		lock.Unlock()
		// Gosched yields the processor, allowing other goroutines to run. It does not
		// suspend the current goroutine, so execution resumes automatically.
		runtime.Gosched()
		if c >= 10 {
			break
		}
	}
	println("----------------------------------------------------------------------------------")
	// NOTE 通道式的信息共享方式
	stringChan := make(chan string)
	producer := func() {
		fmt.Println("send befor ...")
		for i := 1; i <= 10; i++ {
			mesg := "message" + strconv.FormatInt(int64(i), 10)
			stringChan <- mesg

		}
		fmt.Println("produced finished")
	}
	go producer()
	// for {
	// 	consumer := <-stringChan
	// 	fmt.Println(consumer)
	// }
	for consumer := range stringChan {
		fmt.Println(consumer)
		break
		// NOTE 接收方需要接收发送方的一个退出标志，停止接收，否则会一直阻塞，当然发送方在发送停止接收标志后，
		//      也不能在继续发送数据，否则也会阻塞，因为没有goroutine继续接收了。
	}
	println("------ main routine end ------")
}

// -count=1
// -race

// NOTE 竞态测试与原子函数的使用
func TestSharedMemoryCommunicate(t *testing.T) {
	var (
		count    int32
		shutdown int64
		wg       sync.WaitGroup
		// wg1      sync.WaitGroup
		// wg2      sync.WaitGroup
	)
	atomicFunc := func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			atomic.AddInt32(&count, 1)
			// value := count    // read
			runtime.Gosched() // 让当前goroutine暂停,退回执行队列,使其他等待的goroutine运行,让竞争表现明显。
			// value++
			// count = value // write
		}
	}

	wg.Add(2)
	go atomicFunc()
	go atomicFunc()
	wg.Wait()
	fmt.Println(count)

	fmt.Println("--------------------------------------------------------------------")
	// Go语言提供了传统的同步 goroutine 的机制，就是对共享资源加锁。
	// atomic 和 sync 包里的一些函数就可以对共享的资源进行加锁操作。
	// 原子函数能够以很底层的加锁机制来同步访问整型变量和指针
	// LoadInt64 和 StoreInt64。这两个函数提供了一种安全地读和写一个整型值的方式.
	doWork := func(name string) {
		defer wg.Done()
		for {
			fmt.Printf("Doing %s Work\n", name)
			time.Sleep(250 * time.Millisecond)
			if atomic.LoadInt64(&shutdown) == 1 {
				fmt.Printf("Shutting %s Down\n", name)
				break
			}
		}
	}

	wg.Add(2)
	go doWork("A")
	go doWork("B")
	time.Sleep(1 * time.Second)
	fmt.Println("Shutdown Now")
	atomic.StoreInt64(&shutdown, 1)
	wg.Wait()
	// StoreInt64 函数来安全地修改 shutdown 变量的值。如果哪个 doWork goroutine 试图在 main 函数调用
	// StoreInt64 的同时调用 LoadInt64 函数，那么原子函数会将这些调用互相同步，保证这些操作都是安全的，
	// 不会进入竞争状态
	fmt.Println("--------------------------------------------------------------------")
	// A Mutex must not be copied after first use.
	var mutex sync.Mutex // ??? 为什么没有实例化 ？？？
	fmt.Println(&mutex)
	mutexFunc := func(id int) {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			//同一时刻只允许一个goroutine进入这个临界区
			mutex.Lock()
			{
				value := count
				runtime.Gosched()
				// 强制将当前 goroutine 退出当前线程后，调度器会再次分配这个 goroutine 继续运行
				value++
				count = value
			}
			mutex.Unlock() //释放锁，允许其他正在等待的goroutine进入临界区
		}
	}
	wg.Add(2)
	go mutexFunc(1)
	go mutexFunc(2)
	wg.Wait()
	fmt.Println(count)
}

/*
	在 Go语言程序运行时（runtime）实现了一个小型的任务调度器。这套调度器的工作原理类似于操作系统调度线程，
	Go 程序调度器可以高效地将 CPU 资源分配给每一个任务。传统逻辑中，开发者需要维护线程池中线程与 CPU 核心
	数量的对应关系。同样的，Go 地中也可以通过 runtime.GOMAXPROCS() 函数做到.
		runtime.GOMAXPROCS(runtime.NumCPU())
			<1：不修改任何数值。
			=1：单核心执行。
			>1：多核并发执行。
	GOMAXPROCS 同时也是一个环境变量，在应用程序启动前设置环境变量也可以起到相同的作用。k
*/

/*
	TODO Go语言并发并行
	http://c.biancheng.net/view/95.html -> (9.6, 9.7)
*/
