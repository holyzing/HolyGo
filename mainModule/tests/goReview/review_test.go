package goReview

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"testing"
)

/**
三个协程(线程)交替打印自己的编号
*/
func Test3Chan(t *testing.T) {
	chan1, chan2, chan3 := make(chan int), make(chan int), make(chan int)

	gt := func(name string, number int, numChan chan int) {
		for {
			num := <-numChan
			fmt.Println(name, "->", number)
			if num == 1 {
				chan2 <- 2
			} else if num == 2 {
				chan3 <- 3
			} else if num == 3 {
				chan1 <- 1
			}
		}
	}
	go gt("1", 1, chan1)
	go gt("2", 2, chan2)
	go gt("3", 3, chan3)
	chan1 <- 1
	for {
		// 阻塞主协程的三种方法
		// 1- chan
		// 2- WaitGroup
		// 3- 死循环,不停的让出调度
		runtime.Gosched() // 让出并不意味着挂起
	}
}

func TestOptimisticLock(t *testing.T) {
	// lock := &sync.Mutex{}
	nums := [...]int{1, 2, 3}
	currentIndex, currentNum := 0, nums[0]

	printNumer := func(num int) {
		println("starting...", num)
		for {
			// println("running...", num)
			if currentNum == num {
				// lock.Lock()
				fmt.Println(num)
				if currentIndex == len(nums)-1 {
					currentIndex = 0
				} else {
					currentIndex = currentIndex + 1
				}
				currentNum = nums[currentIndex]
				// lock.Unlock()

				// atomic.AddInt64()
				// NOTE 乐观锁,其实都用不着悲观锁
			}
		}
	}
	for _, num := range nums {
		print(num)
		go printNumer(num)
	}
	for {
		runtime.Gosched() // 让出并不意味着挂起
	}
}

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
