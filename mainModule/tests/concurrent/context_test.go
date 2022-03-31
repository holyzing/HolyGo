package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	/*
			type Context interface {
				Deadline() (deadline time.Time, ok bool)
				Done() <-chan struct{}
				Err() error
				Value(key interface{}) interface{}
			}

		    Deadline — 返回 context.Context 被取消的时间，也就是完成工作的截止日期；
			Done : 返回一个 Channel，这个 Channel 会在当前工作完成或者上下文被取消后关闭，多次调用 Done 方法会返回同一个 Channel；
			Err  : 返回 context.Context 结束的原因，它只会在 Done 方法对应的 Channel 关闭时返回非空的值；
				   如果 context.Context 被取消，会返回 Canceled 错误；
				   如果 context.Context 超时，会返回 DeadlineExceeded 错误；
			Value: 从 context.Context 中获取键对应的值，对于同一个上下文来说，多次调用 Value 并传入相同的 Key 会返回相同的结果，
		           该方法可以用来传递请求特定的数据；

			在 Goroutine 构成的树形结构中对信号进行同步以减少计算资源的浪费是 context.Context 的最大作用。
			Go 服务的每一个请求都是通过单独的 Goroutine 处理的，HTTP/RPC 请求的处理器会启动新的 Goroutine 访问数据库和其他服务。

			可能会创建多个 Goroutine 来处理一次请求，而 context.Context 的作用是在不同
			Goroutine 之间同步请求特定数据、取消信号以及处理请求的截止日期。

			每一个 context.Context 都会从最顶层的 Goroutine 一层一层传递到最下层。
			context.Context 可以在上层 Goroutine 执行出现错误时，将信号及时同步给下层。
	*/

	/*
			Context
			emptyCtx 之所以不是一个空的struct{},是因为每一个新建的(顶层)emptyCtx应该都有自己的地址

		    valueCtx

		    cancelCtx
			timerCtx
	*/

	handle := func(ctx context.Context, duration time.Duration, cancelFunc context.CancelFunc) {
		// 这里在开一个协程,用来访问数据啥的,然后主协程阻塞,等待超时或者父协程通知结束

		go func() {
			count := 0
			for {
				if count == 0 {
					fmt.Println("开始工作 ...")
				} else {
					time.Sleep(time.Second)
				}
				count++
				fmt.Println("已经工作了", count, "秒!")
			}
		}()

		select {
		// 其它协程可能会取消 该Context,所以也可以接收取消信号,做后续处理
		case <-ctx.Done():
			fmt.Println("handle", ctx.Err())
		// 超时后取消该Contxt,做后续处理
		case <-time.After(duration):
			fmt.Println("process request with", duration)
			cancelFunc()
		}
	}

	bg := context.Background()                // emptyCtx
	ctx, cancelFunc := context.WithCancel(bg) // timerCtx -> cancelCtx
	defer cancelFunc()
	go handle(ctx, 3*time.Second, cancelFunc)

	// 阻塞主协程,等待子协程结束信号, 这种方式可能会导致其它子协程无法及时处理后续工作,建议使用WaitGroup
	<-ctx.Done()
	fmt.Println("main", ctx.Err())
}

// NOTE context 对于 with_value 是线程安全的, 每次添加键值对都是以当前context作为父节点,衍生一个valueCtx节点,
// 对于 获取键值是从底层开始, 层层向上,直到找到键,如果找不到就到 empty Context 返回 nil
// 这只是说明with_value 是安全的, 但是context是没有提供 del 的
