package goReview

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"testing"
)

// ------------------------------------------------

// Worker  工作队列里的工作者
type Worker interface {
	Work()
}

type holdOnWorker struct {
	trigger  sync.Mutex
	finished sync.Mutex
}

// Ring 环状队列
type Ring struct {
	workers []Worker
	pos     int
}

// Run 开始工作
func (r *Ring) Run() {
	for {
		if r.pos == len(r.workers) {
			r.pos = 0
		}
		r.workers[r.pos].Work()
		r.pos++
	}
}

func (w *holdOnWorker) HoldOn(fn func()) {
	w.trigger.Lock()
	fn()
	w.finished.Unlock()
}

func (w *holdOnWorker) Work() {
	w.trigger.Unlock()
	w.finished.Lock()
}

// ??? CPU 可以并行的去访问一块内存吗 ??? 同一时刻,同一内存只能被一个CPU访问 ??? 内核态的控制访问 ???
// NOTE 先Unlock的锁不执行事务, 不可能同时释放锁,如果不交替执行,则会出现死锁

func TestMutex(t *testing.T) {
	newHoldOnWorker := func() *holdOnWorker {
		w := &holdOnWorker{}
		w.trigger.Lock()
		w.finished.Lock()
		return w
	}

	// fakeThread 假的工作线程
	fakeThread := func(id string, holdOn func(func())) {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for {
			holdOn(func() {
				fmt.Println(id)
			})
		}
	}

	w1 := newHoldOnWorker()
	go fakeThread("1", w1.HoldOn)

	w2 := newHoldOnWorker()
	go fakeThread("2", w2.HoldOn)

	w3 := newHoldOnWorker()
	go fakeThread("3", w3.HoldOn)

	r := &Ring{
		workers: []Worker{w1, w2, w3},
	}
	r.Run()

	fmt.Println("r.Run 不是死循环 ???? 不应该是 go r.Run")

	var exit = make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
}

// 锁的原理,锁被创建的时候是不被标记占用的,一旦某个协程占用,那么再次调用这个锁标记的时候,因为已经被标记了,
// 所以会阻塞当前执行流,直到这个锁被清楚标记,然后再次标记

// 实现一个乐观锁，你给123加一个版本号
// atomic add 配上 sleep食用更佳
// TODO 考虑控制反转
