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
		n := rand.Intn(100)
		if n%13 == 0 {
			fmt.Printf("Player %s Missed\n", name)
			// 关闭通道，表示我们输了
			close(court)
			return
		}
		// 显示击球数，并将击球数加1
		fmt.Printf("Player %s Hit %d\n", name, ball)
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
	court := make(chan int)
	// 计数加 2，表示要等待两个goroutine
	wg.Add(2)
	// 启动两个选手
	go Player("Nadal", court)
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
		for {
			select {
			case num := <-ch: // 如果ch成功读到数据，则进行该case处理语句
				fmt.Println("num = ", num)
			case <-time.After(3 * time.Second):
				fmt.Println("超时")
				quit <- true
			}
		}
	}() //别忘了()
	for i := 0; i < 5; i++ {
		ch <- i
		time.Sleep(time.Second)
	}
	<-quit
	fmt.Println("程序结束")
}
