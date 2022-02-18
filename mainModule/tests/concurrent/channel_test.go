package concurrent

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// TODO 通道测试
func TestChannel(t *testing.T) {
	// ??? 通道类型的变量本身就是指针 ？？？，就是说它是引用传递 ？？
	// ??? goroutine 不会先与 mainroutine 结束

	// 把数据往通道中发送时，如果接收方一直都没有接收，那么发送操作将持续阻塞。
	// Go 程序运行时能智能地发现一些永远无法发送成功的语句并做出提示.

	// ① 通道的收发操作在不同的两个 goroutine 间进行。
	// 由于通道的数据在没有接收方处理时，数据发送方会持续阻塞，因此通道的接收必定在另外一个 goroutine 中进行。

	// ② 接收将持续阻塞直到发送方发送数据。
	// 如果接收方接收时，通道中没有发送方发送数据，接收方也会发生阻塞，直到发送方发送数据为止。

	// ③ 通道一次只能接收一个数据元素, 每一个入队元素只能按顺序被读取一次。
	// 在任何时候，同时只能有一个 goroutine 访问通道进行发送和获取数据。换句话说通道的读写是同步的。

	// data := <-ch			// 执行该语句时将会阻塞，直到接收到数据并赋值给 data 变量
	// data, ok := <-ch     // 未接收到数据时，data 为通道类型的零值。ok：表示是否接收到数据。
	// <-ch                 // 放弃接收，主要是为了疏通阻塞
	// 非阻塞的通道接收方法可能造成高的 CPU 占用，因此使用非常少。
	// 要实现接收超时检测，可以配合 select 和计时器 channel 进行

	// NOTE 单向 channel，所谓单向channel 指的是通过一定语法规则，将可写可读的 channel 赋值给
	//      只读或者只写的channel 变量里，channel本身是可写可读的,否则只读和只写是没有什么意义的。
	ch := make(chan int)
	var readOnly <-chan int = ch
	var writeOnly chan<- int = ch

	// readCh := make(<-chan int) // 没什么意义, 注意 （<-chan int）这是一个类型

	// NOTE 定时器是一个只能读取的 chan
	// type Timer struct {
	//	C <-chan Time
	//	r runtimeTimer
	// }

	// NOTE 关闭 channel
	// close(readOnly) // cannot close receive-only channel
	close(writeOnly)

	// NOTE 判断 channel 的状态是否为关闭
	_, ok := <-readOnly
	if ok {
		close(ch)
	} else {
		fmt.Println("chan has been closed")
	}
}

var wg sync.WaitGroup

func Runner(baton chan int) {
	var newRunner int

	runner := <-baton
	fmt.Printf("Runner %d Running With Baton\n", runner)
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d To The Line\n", newRunner)

		go Runner(baton)
		// 再开一个子协程，子线程到 runner := <-baton 阻塞，知道主线程 baton <- newRunner 到达终点
	}

	time.Sleep(100 * time.Millisecond)
	if runner == 4 {
		fmt.Printf("Runner %d Finished, Race Over\n", runner)
		wg.Done()
		return
	}
	fmt.Printf("Runner %d Exchange With Runner %d\n", runner, newRunner)
	baton <- newRunner
}

func Player(name string, court chan int) {
	defer wg.Done()

	for {
		ball, ok := <-court
		if !ok {
			// 如果通道被关闭，我们就赢了
			fmt.Printf("Player %s Won\n", name)
			return
		}
		// 选随机数，然后用这个数来判断我们是否丢球
		n := rand.Intn(100) // TODO 每次运行生成的随机数是一样的
		if n%13 == 0 {
			fmt.Printf("Player %s Missed value %d\n", name, n)
			// 关闭通道，表示我们输了
			close(court)
			return
		}
		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s Hit %d value %d\n", name, ball, n)
		ball++
		// 将球打向对手
		court <- ball
	}
}

func TestUnBufferedChannelWithVR(t *testing.T) {
	/*
		无缓冲的通道（unbuffered channel）是指在接收前没有能力保存任何值的通道。
		这种类型的通道要求发送 goroutine 和接收 goroutine 同时准备好，才能完成发送和接收操作。
		如果两个 goroutine 没有同时准备好，通道会导致先执行发送或接收操作的 goroutine 阻塞等待。
		这种对通道进行发送和接收的交互行为本身就是同步的。其中任意一个操作都无法离开另一个操作单独存在。

		阻塞: 指的是由于某种原因数据没有到达，当前协程（线程）持续处于等待状态，直到条件满足才解除阻塞。
		同步: 指的是在两个或多个协程（线程）之间，保持数据内容一致性的机制。

		在通道内传递数据的过程是不会交出线程占用的。
	*/

	/*
		court := make(chan int)
		// 计数加 2，表示要等待两个goroutine
		wg.Add(2)
		// 启动两个选手
		go Player("Nadal", court)
		// time.Sleep(5*time.Second)
		go Player("Djokovic", court)
		// 发球
		court <- 1
		// 等待游戏结束
		wg.Wait()

		println("-------------------------------------------------------------------------------------")

		// 使用无缓存的通道模拟 接力比赛
		wg.Add(1)
		baton := make(chan int)
		go Runner(baton)
		baton <- 1
		wg.Wait()

		println("-------------------------------------------------------------------------------------")
	*/
	ch := make(chan int)
	wg.Add(10)
	for i := range "0123456789" {
		i := i
		go func() {
			// Loop variables captured by 'func' literals in 'go' statements might have unexpected values
			println(i)
			wg.Done()
			ch <- i
		}()
	}
	wg.Wait()
	println("*****************************************************")
	for i := range ch {
		println(i)
	}
}

// 缓存式的通道
func TestBufferedChannel(t *testing.T) {
	/**
	Go语言中有缓冲的通道（buffered channel）是一种在被接收前能存储一个或者多个值的通道。
	这种类型的通道并不强制要求 goroutine 之间必须同时完成发送和接收。通道会阻塞发送和接收动作的条件也会不同。
	只有在通道中没有要接收的值时，接收动作才会阻塞。只有在通道没有可用缓冲区容纳被发送的值时，发送动作才会阻塞。

	这导致有缓冲的通道和无缓冲的通道之间的一个很大的不同：
	无缓冲的通道保证进行发送和接收的 goroutine 会在同一时间进行数据交换；有缓冲的通道没有这种保证。

	在无缓冲通道的基础上，为通道增加一个有限大小的存储空间形成带缓冲通道。
	带缓冲通道在发送时无需等待接收方接收即可完成发送过程，并且不会发生阻塞，只有当存储空间满时才会发生阻塞。
	同理，如果缓冲通道中有数据，接收时将不会发生阻塞，直到通道中没有数据可读时，通道将会再度阻塞。

	无缓冲通道保证收发过程同步。无缓冲收发过程类似于快递员给你电话让你下楼取快递，整个递交快递的过程是同步发生的，
	你和快递员不见不散。但这样做快递员就必须等待所有人下楼完成操作后才能完成所有投递工作。
	如果快递员将快递放入快递柜中，并通知用户来取，快递员和用户就成了异步收发过程，效率可以有明显的提升。
	带缓冲的通道就是这样的一个“快递柜”。
	*/
	ch := make(chan int, 3)
	fmt.Println(len(ch))
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(len(ch))

	// 带缓冲通道在很多特性上和无缓冲通道是类似的。无缓冲通道可以看作是长度永远为 0 的带缓冲通道。
	// 因此根据这个特性，带缓冲通道在下面列举的情况下依然会发生阻塞：
	// 带缓冲通道被填满时，尝试再次发送数据时发生阻塞。
	// 带缓冲通道为空时，尝试接收数据时发生阻塞。

	// 为什么Go语言对通道要限制长度而不提供无限长度的通道？
	// 我们知道通道（channel）是在两个 goroutine 间通信的桥梁。
	// 使用 goroutine 的代码必然有一方提供数据，一方消费数据。
	// 当提供数据一方的数据供给速度大于消费方的数据处理速度时，如果通道不限制长度，那么内存将不断膨胀直到应用崩溃。
	// so,限制通道的长度利于约束数据提供方的供给速度,供给数据量必须在消费方处理量+通道长度的范围内,才能正常地处理数据。
}

// 操作已经关闭的channel
func TestOperateClosedChannel(t *testing.T) {
	unbufferedChannel := make(chan int)
	fmt.Println("unbufferedChannel:", cap(unbufferedChannel), len(unbufferedChannel))
	close(unbufferedChannel)
	fmt.Println("unbufferedChannel:", cap(unbufferedChannel), len(unbufferedChannel))
	res := <-unbufferedChannel
	fmt.Println(res)
	// NOTE：1. 非缓存通道关闭，不能发送数据
	// unbufferedChannel <- 1  panic: send on closed channel
	// NOTE：2. 非缓存通道关闭，可以接收数据，接收的数据是通道数据类型的零值,且会立即返回

	// 非缓存chan的长度和容量始终为0,是因为在读取方准备好之后,运行时会将发送方的数据直接拷贝到接收方.
	// 无论是写接收还是写发送都是可以的,但是对于非缓冲chan来说两种操作是强依赖的,都会阻塞所在协程.
	println("---------------------------------------------------------------------")

	bufferedChannel := make(chan int, 3)
	fmt.Println("bufferedChannel:", cap(bufferedChannel), len(bufferedChannel))
	bufferedChannel <- 1 // 数据满则阻塞
	fmt.Println("bufferedChannel:", cap(bufferedChannel), len(bufferedChannel))
	e := <-bufferedChannel // 无数据会阻塞
	fmt.Println("bufferedChannel:", e, cap(bufferedChannel), len(bufferedChannel))
	// NOTE：3. 缓存通道关闭，不能发送数据
	// bufferedChannel <- 1
	// NOTE：4. 缓存通道关闭，可以接收数据，接收的数据是通道数据缓存的值，如果无缓存则是通道类型的零值
	// d <- bufferedChannel // (send to non-chan type int)
	close(bufferedChannel)
	fmt.Println("bufferedChannel:", cap(bufferedChannel), len(bufferedChannel))
	e = <-bufferedChannel
	fmt.Println(e)
	println("---------------------------------------------------------------------")

	/*
		通道是一个引用对象，和 map 类似。map 在没有任何外部引用时，Go语言程序在运行时（runtime）会自动对内
		存进行垃圾回收（Garbage Collection, GC）。类似的，通道也可以被垃圾回收，但是通道也可以被主动关闭。

		从已经关闭的通道接收数据或者正在接收数据时，将会接收到通道类型的零值，然后停止阻塞并返回。
	*/

	// 创建一个整型带两个缓冲的通道
	ch := make(chan int, 2)

	// 给通道放入两个数据
	ch <- 0
	ch <- 1

	// 关闭通道。此时，带缓冲通道的数据不会被释放，通道也没有消失
	close(ch)
	for i := 0; i <= cap(ch); i++ {

		// 缓冲通道在关闭后依然可以访问内部的数据。
		v, ok := <-ch

		// 如果越界访问，则访问返回的是通道类型的零值
		fmt.Println(v, ok)
	}
}

/**
golang中分为值类型和引用类型
	值类型分别有：int系列、float系列、bool、string、数组和结构体
	引用类型有：指针、slice切片、管道channel、接口interface、map、函数等
	值类型的特点是：变量直接存储值，内存通常在栈中分配
	引用类型的特点是：变量存储的是一个地址，这个地址对应的空间里才是真正存储的值，内存通常在堆中分配
*/
