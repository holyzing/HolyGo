package concurrent

import (
	"fmt"
	"testing"
)

func TestReadChan(t *testing.T) {
	nB := make(chan string, 5)
	go func() {
		nB <- "a"
		nB <- "b"
		nB <- "c"
		// close(nB)
	}()
	v, ok := <-nB
	fmt.Println(v, ok)
	v, ok = <-nB
	fmt.Println(v, ok)
	v, ok = <-nB
	fmt.Println(v, ok)

	v, ok = <-nB
	fmt.Println("-------------->", v, ok)

	v, ok = <-nB
	fmt.Println(v, ok)
}

// 1. 写未关闭：写入或者阻塞
// 2. 写已关闭：panic
// 3. 读未关闭：读出或者阻塞
// 4. 读已关闭：反复可读，有值则为True，无值则为False
