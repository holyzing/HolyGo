package goReview

import (
	"log"
	"runtime"
	"testing"
	"time"
)

/**
三个协程(线程)交替打印自己的编号
*/
func Test3Chan(t *testing.T) {
	chan1, chan2, chan3 := make(chan int), make(chan int), make(chan int)

	gt := func(name string, number int, numChan chan int) {
		for {
			num := <-numChan
			log.Println(name, "->", number)
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

func TestCloseChan(t *testing.T) {
	c := make(chan int)

	go func(cc chan int) {
		log.Println("start ....")
		log.Println(<-cc)
		log.Println("end ....")
		select {}
		// TODO 如何合理的 Close Chan,且一定要记得关闭Chan，负责会造成内存泄露
		// NOTE 所有协程阻塞，则会Panic，Select 是阻塞操作
		// NODE 关闭的Channel 不能写入，但是依旧可以读，读出的是零值
	}(c)

	time.Sleep(5 * time.Second)
	close(c)
	time.Sleep(5 * time.Second)
	log.Println(<-c)

	select {}
}
