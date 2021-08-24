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
}

func TestTime(t *testing.T) {
	var tm = time.Now()
	fmt.Println(tm)
	fmt.Println(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond())
	fmt.Println(tm.Location())
	// fmt.Println(time.Date(2021, 12, 17, 20, 50, 50, 500))
}
