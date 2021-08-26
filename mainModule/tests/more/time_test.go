package more

import (
	"fmt"
	"mainModule/collects"
	"mainModule/tests"
	"testing"
	"time"
)

// NOTE GO语言中的访问控制
// NOTE 所有语法可见性均定义在package这个级别
func TestPrivateOrPublic(t *testing.T) {
	collects.Cao()
	// NOTE 只能访问到结构体的公有的属性
	se := tests.Second{}
	fmt.Println(se.A)
	m := More{A: 12, b: false}
	println(m.A, m.b)
}

func TestTime(t *testing.T) {
	// time.Time
	// time.Duration
	// time.Location
	// time.Timer
	// time.Ticker

	var tm = time.Now()
	fmt.Println(tm)
	tm, err := time.Parse("2006-01-02 15:04:05", "2011-04-23 12:24:51")
	if err != nil {
		fmt.Println(tm)
	}
	fmt.Println(time.Local)

	fmt.Println(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond())
	fmt.Println(tm.Location())
	// fmt.Println(time.Date(2021, 12, 17, 20, 50, 50, 500))
}
