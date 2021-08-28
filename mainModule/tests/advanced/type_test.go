package advanced

import (
	"fmt"
	"reflect"
	"testing"
)

// NOTE 定义结构体
type Mutex struct{}

// NOTE 定义接口
type Animal interface {
}

// NOTE 定义新类型
type MyInt int

func (m MyInt) ShowValue() {
	fmt.Println("show int", m)
}

// NOTE 定义类型别名
type B32 = int32

type nameMap = map[string]interface{}

type FuncType func(s string)

func (ft FuncType) ServerCb() error {
	ft("cbb")
	fmt.Println("server cb")
	return nil
}

// NOTE var a interface{} = nil

func TestType(t *testing.T) {
	f := func(s string) {
		fmt.Println("sss", s)
	}
	ft := FuncType(f) // Type Conversion
	ft.ServerCb()

	defaultInt := 12
	fmt.Println(reflect.TypeOf(defaultInt))
	convertedInt := int64(defaultInt)
	fmt.Println(reflect.TypeOf(convertedInt))

	var AnyType interface{}
	AnyType = 1
	AnyType = 'a'
	AnyType = "dasds"
	AnyType = [...]int{4, 6: 6}
	fmt.Println(AnyType)
}
