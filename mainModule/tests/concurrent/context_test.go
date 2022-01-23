package concurrent

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T)  {
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


	handle := func (ctx context.Context, duration time.Duration) {
		select {
		case <-ctx.Done():
			fmt.Println("handle", ctx.Err())
		case <-time.After(duration):
			fmt.Println("process request with", duration)
		}
	}

	bg := context.Background()  // emptyCtx
	ctx, cancelFunc := context.WithTimeout(bg, 1*time.Second)  // timerCtx -> cancelCtx
	defer cancelFunc()
	go handle(ctx, 500*time.Millisecond)

	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
}