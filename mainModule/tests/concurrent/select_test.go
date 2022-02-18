package concurrent

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

// 使用 select解决超时问题
func TestSelectTimeout(t *testing.T) {
	/*
		虽然 select 机制不是专门为超时而设计的，却能很方便的解决超时问题，
		因为 select 的特点是只要其中有一个 case 已经完成，程序就会继续往下执行，而不会考虑其他 case 的情况。

		超时机制本身虽然也会带来一些问题，比如在运行比较快的机器或者高速的网络上运行正常的程序，到了慢速的机器或者
		网络上运行就会出问题，从而出现结果不一致的现象，但从根本上来说，解决死锁问题的价值要远大于所带来的问题。

		与 switch 语句相比，select 有比较多的限制，其中最大的一条限制就是每个 case 语句里必须是一个 IO 操作.

		select {
			case <-chan1:
			// 如果chan1成功读到数据，则进行该case处理语句
			case chan2 <- 1:
			// 如果成功向chan2写入数据，则进行该case处理语句
			default:
			// 如果上面都没有成功，则进入default处理流程
		}

		在一个 select 语句中，Go语言会按顺序从头至尾评估每一个发送和接收的语句。
		如果其中的任意一语句可以继续执行（即没有被阻塞），那么就从那些可以执行的语句中任意选择一条来使用。
		如果没有任意一条语句可以执行（即所有的通道都被阻塞），那么有如下两种可能的情况：
			如果给出了 default 语句，那么就会执行 default 语句，同时程序的执行会从 select 语句后的语句中恢复；
			如果没有 default 语句，那么 select 语句将被阻塞，直到至少有一个通信可以进行下去。
	*/
	ch := make(chan int)
	quit := make(chan bool)
	//新开一个协程
	go func() {
	outer:
		for {
			select {
			// 如果ch成功读到数据，则进行该case处理语句
			case num := <-ch:
				fmt.Println("num = ", num)
			case <-time.After(3 * time.Second):
				fmt.Println("超时")
				quit <- true
				break outer
			}
		}
	}()
	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	<-quit
	fmt.Println("程序结束")
	println("-----------------------------------------------------")

	/*  NOTE Select 多路复用
	在使用通道时，想同时接收多个通道的数据是一件困难的事情。通道在接收数据时，
	如果没有数据可以接收将会发生阻塞。虽然可以使用for 循环依次进行遍历，但运行性能会非常差。
	*/

	for {
		// ??? 如果 两个 case 都同时有IO 触发，会都执行吗 ？
		select {
		case ch <- 0:
		case ch <- 1:
		}
		i := <-ch
		fmt.Println("Value received:", i)
		// 随机读写0或者1
	}
}

// 使用select 模拟远程过程调用
func TestSelectRpc(t *testing.T) {
	/*
		服务器开发中会使用RPC（Remote Procedure Call，远程过程调用）简化进程间通信的过程。
		RPC 能有效地封装通信过程，让远程的数据收发通信过程看起来就像本地的函数调用一样。
		使用通道代替 Socket 实现 RPC 的过程。
		客户端与服务器运行在同一个进程，服务器和客户端在两个 goroutine 中运行
	*/
	RPCClient := func(server chan string, req string) (string, error) {
		server <- req
		var err error
		select {
		case resp, ok := <-server:
			if ok {
				return resp, nil
			}
		case currentTime := <-time.After(3):
			fmt.Println(currentTime)
			err = errors.New("time out !")
		}
		return "", err
	}

	RPCServer := func(req chan string) {
		// 模拟超时 time.Sleep(4 * time.Second)
		for {
			// 接收客户端请求
			data := <-req
			// 打印接收到的数据
			fmt.Println("server received:", data)
			// 反馈给客户端收到
			req <- "roger"
		}
	}

	ch := make(chan string)
	// 并发执行服务器逻辑
	go RPCServer(ch)
	// 客户端请求数据和接收数据
	recv, err := RPCClient(ch, "hi")
	if err != nil {
		// 发生错误打印
		fmt.Println(err)
	} else {
		// 正常接收到数据
		fmt.Println("client received", recv)
	}

}

// 通道响应计时器的事件
func TestNewTimerNewTickerAfter(t *testing.T) {
	// Go语言中的通道和 goroutine 的设计，定时任务可以在 goroutine
	// 中通过同步的方式完成，也可以通过在 goroutine 中异步回调完成

	exit := make(chan int)
	callback := func() {
		fmt.Println("callback ...")
		exit <- 1
	}
	timer := time.AfterFunc(3*time.Second, callback)
	// timer.Stop()
	fmt.Println(timer)
	<-exit
	// NOTE 子协程不会像java或者python中的线程可以 join 到 主协程

	/*
		time.After() 函数是在 time.NewTimer() 函数上进行的封装，timer.NewTimer() 和 time.NewTicker()。
		计时器（Timer）的原理和倒计时闹钟类似，都是给定多少时间后触发。
		打点器（Ticker）的原理和钟表类似，钟表每到整点就会触发。
		这两种方法创建后会返回 time.Ticker 对象和 time.Timer 对象，里面通过一个 C 成员，
		类型是只能接收的时间通道（<-chan Time），使用这个通道就可以获得时间触发的通知。
	*/

	// --------------------------------------------------

	// 创建一个打点器, 每500毫秒触发一次
	ticker := time.NewTicker(time.Millisecond * 500)
	// 创建一个计时器, 2秒后触发
	stopper := time.NewTimer(time.Second * 2)
	// 声明计数变量
	var i int
	// 不断地检查通道情况
	for {
		// 多路复用通道
		select {
		case <-stopper.C: // 计时器到时了
			fmt.Println("stop")
			// 跳出循环
			goto StopHere
		case <-ticker.C: // 打点器触发了
			// 记录触发了多少次
			i++
			fmt.Println("tick", i)
		}
	}
	// 退出的标签, 使用goto跳转
StopHere:
	fmt.Println("done")
}

// select 模拟生产者消费者模型
func TestPCModel(t *testing.T) {
	pump1 := func(ch chan int) {
		for i := 0; ; i++ {
			ch <- i * 2
		}
	}

	pump2 := func(ch chan int) {
		for i := 0; ; i++ {
			ch <- i + 5
		}
	}

	suck := func(ch1, ch2 chan int) {
	forTag:
		for {
			select {
			case v := <-ch1:
				// The 'fallthrough' statement is out of place
				// fallthrough
				fmt.Printf("Received on channel 1: %d\n", v)
			case v := <-ch2:
				fmt.Printf("Received on channel 2: %d\n", v)
				// 不加 fortag break 的是select 而不是 for
				break forTag
			}
		}
	}

	ch1 := make(chan int)
	ch2 := make(chan int)
	go pump1(ch1)
	go pump2(ch2)
	go suck(ch1, ch2)
	time.Sleep(1e9)
}
