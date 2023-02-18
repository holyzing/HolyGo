package goReview

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

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
				time.Sleep(1 * time.Second)
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
