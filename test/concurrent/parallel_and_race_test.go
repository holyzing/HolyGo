package concurrent

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelWithMultiCore(t *testing.T) {
	// 服务器的处理器大都是单核频率较低而核心数较多，对于支持高并发的程序语言，
	// 可以充分利用服务器的多核优势，从而降低单核压力，减少性能浪费。
	AsyncFunc := func(index int) {
		sum := 0
		for i := 0; i < 10000; i++ {
			sum += 1
		}
		fmt.Printf("线程%d, sum为:%d\n", index, sum)
	}
	for i := 0; i < 5; i++ {
		go AsyncFunc(i)
	}
	time.Sleep(10 * time.Second)
	// 主协程很快结束

	// 在执行一些昂贵的计算任务时，希望能够尽量利用现代服务器普遍具备的多核特性来尽量将任务并行化，从而达到降低总计算
	// 时间的目的。此时需要了解 CPU 核心的数量，并针对性地分解计算任务到多个 goroutine 中去并行运行。

	/**
	type Vector []float64

	var NCPU = runtime.NumCPU() // 假设总共有16核

	func (v Vector) DoSome(i, n int, u Vector, c chan int) {
		for ; i < n; i++ {
			v[i] += u.Op(v[i])
		}
		// 发信号告诉任务管理者已经计算完成
		c <- 1
	}

	func (v Vector) DoAll(u Vector) {
		// 用于接收每个CPU的任务完成信号
		c := make(chan int, NCPU)
		for i := 0; i < NCPU; i++ {
			go v.DoSome(i*len(v)/NCPU, (i+1)*len(v)/NCPU, u, c)
		}
		// 等待所有CPU的任务完成
		for i := 0; i < NCPU; i++ {
			// 获取到一个数据，表示一个CPU计算完成了
			<-c
		}
		// 到这里表示所有计算已经结束
	}
	*/

	// 模拟一个完全可以并行的计算任务：计算 N 个整型数的总和。可以将所有整型数分成 M 份，M 即 CPU 的个数。
	// 让每个 CPU 开始计算分给它的那份计算任务，最后将每个 CPU 的计算结果再做一次累加，得到所有 N 个整型数的总和。

	/*
		是否可以将总的计算时间降到接近原来的 1/N 呢？答案是不一定。如果掐秒表，会发现总的执行时间没有明显缩短。
		再去观察 CPU 运行状态，发现尽管有 16 个 CPU 核心，但在计算过程中其实只有一个 CPU 核心处于繁忙状态，
		这是会让很多Go语言初学者迷惑的问题。

		官方给出的答案是，这是当前版本（Go1.13.4）的 Go 编译器还不能很智能地去发现和利用多核的优势。
		虽然确实创建了多个 goroutine，并且从运行状态看这些 goroutine 也都在并行运行，但实际上所有这些 goroutine
		都运行在同一个 CPU 核心上，在一个 goroutine 得到时间片执行的时候，其他 goroutine 都会处于等待状态。
		从这一点可以看出，虽然 goroutine 简化了写并行代码的过程，但实际上整体运行效率并不真正高于单线程程序。

		虽然Go语言还不能很好的利用多核心的优势，但也可以先通过设置环境变量 GOMAXPROCS 的值来控制使用多少个CPU核心。
		通过直接设置环境变量 GOMAXPROCS 的值，或者在代码中启动 goroutine 之前先调用
		runtime.GOMAXPROCS(runtime.NumCPU()) 设置使用所有核心：
	*/
}

func TestTelnetSimulate(t *testing.T) {
	/**
	Telnet 协议是 TCP/IP 协议族中的一种。它允许用户（Telnet 客户端）通过一个协商过程与一个远程设备进行通信。
	在操作系统中可以在命令行使用 Telnet 命令发起 TCP 连接。
	我们一般用 Telnet 来连接 TCP 服务器，键盘输入一行字符回车后，即被发送到服务器上。

	本例将使用一部分 Telnet 协议与服务器进行通信。
	服务器的网络库为了完整展示自己的代码实现了完整的收发过程，一般比较倾向于使用发送任意封包返回原数据的逻辑。
	这个过程类似于对着大山高喊，大山把你的声音原样返回的过程。也就是回音（Echo）。
	使用 Go语言中的 Socket、goroutine 和通道编写一个简单的 Telnet 协议的回音服务器。

	回音服务器的代码分为 4 个部分，分别是接受连接、会话处理、Telnet 命令处理和程序入口。
	*/

	println("-------------------------------------------------------------------------------------")

	processTelnetCommand := func(str string, exitChan chan int) bool {
		// @close指令表示终止本次会话
		if strings.HasPrefix(str, "@close") {
			fmt.Println("Session closed")
			// 告诉外部需要断开连接
			return false
			// @shutdown指令表示终止服务进程
		} else if strings.HasPrefix(str, "@shutdown") {
			fmt.Println("Server shutdown")
			// 往通道中写入0, 阻塞等待接收方处理
			exitChan <- 0
			// 告诉外部需要断开连接
			return false
		}
		// 打印输入的字符串
		fmt.Println(str)
		return true
	}

	// 每个连接的会话就是一个接收数据的循环。当没有数据时，调用 reader.ReadString 会发生阻塞，等待数据的到来。
	// 一旦数据到来，就可以进行各种逻辑处理。
	// 回音服务器的基本逻辑是“收到什么返回什么”，reader.ReadString 可以一直读取 Socket
	// 连接中的数据直到碰到期望的结尾符。这种期望的结尾符也叫定界符，一般用于将 TCP 封包中的逻辑数据拆分开。
	// 下例中使用的定界符是回车换行符（“\r\n”），HTTP 协议也是使用同样的定界符。
	// 使用 reader.ReadString() 函数可以将封包简单地拆分开。

	handleSession := func(conn net.Conn, exitChan chan int) {
		fmt.Println("Session started:")
		// 创建一个网络连接数据的读取器
		reader := bufio.NewReader(conn)
		// 接收数据的循环
		for {
			// 读取字符串, 直到碰到回车返回 内部会自动处理粘包过程，直到下一个回车符到达后返回数据。
			str, err := reader.ReadString('\n')
			// 数据读取正确
			if err == nil {
				// 去掉字符串尾部的回车
				str = strings.TrimSpace(str)
				// 处理Telnet指令
				if !processTelnetCommand(str, exitChan) {
					conn.Close()
					break
				}
				// Echo逻辑, 发什么数据, 原样返回
				conn.Write([]byte(str + "\r\n"))
			} else {
				// 发生错误
				fmt.Println("Session closed")
				conn.Close()
				break
			}
		}
	}

	// 回音服务器能同时服务于多个连接。要接受连接就需要先创建侦听器，侦听器需要一个侦听地址和协议类型。
	// 一个服务器可以开启多个侦听器
	// 主机 IP：一般为一个 IP 地址或者域名，127.0.0.1 表示本机地址。
	// 端口号：16 位无符号整型值，一共有 65536 个有效端口号。
	// 在会话中处理的操作和接受连接的业务并不冲突可以同时进行,使用 goroutine 轻松实现会话处理和接受连接的并发执行.

	Server := func(address string, exitChan chan int) {
		// 根据给定地址进行侦听
		l, err := net.Listen("tcp", address)
		// 如果侦听发生错误, 打印错误并退出
		if err != nil {
			fmt.Println(err.Error())
			// 往 exitChan 写入一个整型值时，进程将以整型值作为程序返回值来结束服务器。
			exitChan <- 1
		}
		// 打印侦听地址, 表示侦听成功
		fmt.Println("listen: " + address)
		// 延迟关闭侦听器
		defer l.Close()
		// 侦听循环
		for {
			// 新连接没有到来时, Accept是阻塞的
			conn, err := l.Accept()
			// *tcp.Conn
			// 发生任何的侦听错误, 打印错误并退出服务器
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			// 根据连接开启会话, 这个过程需要并行执行
			go handleSession(conn, exitChan)
		}
	}

	// 创建一个程序结束码的通道
	exitChan := make(chan int)
	// 将服务器并发运行
	go Server("127.0.0.1:7001", exitChan)
	// 通道阻塞, 等待接收返回值
	code := <-exitChan
	// 标记程序返回值并退出
	os.Exit(code)

}

// -------------------------------------------------------------------------------------------------

// Go语言程序可以使用通道进行多个 goroutine 间的数据交换，但这仅仅是数据同步中的一种方法。
// 通道内部的实现依然使用了各种锁，因此优雅代码的代价是性能。
// 在某些轻量级的场合，原子访问（atomic包）、互斥锁（sync.Mutex）以及等待组（sync.WaitGroup）能最大程度满足需求。

// seq := 2  // expected declaration

func TestAtomic(t *testing.T) {
	var seq int64

	// 本例中只是对变量进行增减操作，虽然可以使用互斥锁（sync.Mutex）解决竞态问题，但是对性能消耗较大。
	// 在这种情况下，推荐使用原子操作（atomic）进行变量操作。

	// 序列号生成器
	GenID := func() int64 {
		// 尝试原子的增加序列号
		atomic.AddInt64(&seq, 1)
		// 没有使用 atomic.AddInt64() 的返回值作为 GenID() 函数的返回值，因此会造成一个竞态问题。
		// 操作同一 seq 存在竞态
		return seq
		// ??? return 是原子的吗 ？？？？
	}
	//生成10个并发序列号
	for i := 0; i < 10; i++ {
		go GenID()
	}
	fmt.Println(GenID())
}

// 互斥锁
func TestRWMutex(t *testing.T) {
	var (
		// 逻辑中使用的某个变量
		count int
		// 与变量对应的使用互斥锁
		countGuard sync.RWMutex
	)

	GetCount := func() int {
		// 锁定
		// 一般情况下，建议将互斥锁的粒度设置得越小越好，降低因为共享访问时等待的时间
		countGuard.Lock()
		// 在函数退出时解除锁定
		defer countGuard.Unlock()
		return count
	}

	SetCount := func(c int) {
		// 一旦发生加锁，如果另外一个 goroutine 尝试继续加锁时将会发生阻塞，直到这个countGuard被解锁。
		countGuard.Lock()
		count = c
		countGuard.Unlock()
	}
	// 可以进行并发安全的设置
	SetCount(1)
	// 可以进行并发安全的获取
	fmt.Println(GetCount())

	// ---------------------------------------------------------------------------------------------
	// 经典的单写多读模型。在读锁占用的情况下，会阻止写，但不阻止读，也就是多个 goroutine 可同时获取读锁
	// （调用 RLock() 方法；而写锁（调用 Lock() 方法）会阻止任何其他 goroutine（无论读和写）进来，
	// 整个锁相当于由该 goroutine 独占。从 RWMutex 的实现看，RWMutex 类型其实组合了 Mutex：
	// 在读多写少的环境中，可以优先使用读写互斥锁（sync.RWMutex），它比互斥锁更加高效。

}

func TestWaitGroup(t *testing.T) {
	// Go语言中除了可以使用通道（channel）和互斥锁进行两个并发程序间的同步外，
	// 还可以使用等待组进行多个任务的同步，等待组可以保证在并发环境中完成指定数量的任务.
	// 在 sync.WaitGroup（等待组）类型中，每个 sync.WaitGroup 值在内部维护着一个计数，此计数的初始默认值为零。

	// 声明一个等待组
	var wg sync.WaitGroup
	// 准备一系列的网站地址
	var urls = []string{
		"http://www.github.com/",
		"https://www.qiniu.com/",
		"https://www.golangtc.com/",
	}
	// 遍历这些地址
	for _, url := range urls {
		// 每一个任务开始时, 将等待组增加1
		wg.Add(1)
		// 开启一个并发
		go func(url string) {
			// 使用defer, 表示函数完成时将等待组值减1
			defer wg.Done() // 等价于 wg.add(-1)
			// 使用http访问提供的地址
			_, err := http.Get(url)
			// 访问完成后, 打印地址和可能发生的错误
			fmt.Println(url, err)
			// 通过参数传递url地址
		}(url)
	}
	// 等待所有的任务完成
	wg.Wait()
	fmt.Println("over")
	/*
		当一个协程调用了 wg.Wait() 时，如果此时 wg 维护的计数为零，则此 wg.Wait() 此操作为一个空操作（noop）；
		否则（计数为一个正整数），此协程将进入阻塞状态。当以后其它某个协程将此计数更改至 0 时（一般通过调用 wg.Done()），
		此协程将重新进入运行状态（即 wg.Wait() 将返回）。
	*/
}
