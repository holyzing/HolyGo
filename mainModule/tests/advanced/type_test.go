package advanced

import (
	"fmt"
	"reflect"
	"testing"
	"time"
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

	// 顶层类型
	var AnyType interface{}
	AnyType = 1
	AnyType = 'a'
	AnyType = "dasds"
	AnyType = [...]int{4, 6: 6}
	fmt.Println(AnyType)
}

// TODO rune iota int

// TODO 类型断言
func TestTypeAssert(t *testing.T) {
	var i interface{} = 1
	i2 := i.(int) // WHAT 类型转换 ???
	print(i2)

	var i8 int8 = 63
	print(int(i8))

	// 1630001-01-01 00:00:00 +0000 UTC
	var ti time.Time
	fmt.Println(ti)
}

// -------------------------------------------------------------------------------------------------

type I interface {
	show(string)
	say(s int)
}

type S1 struct {
	age int
}

type S2 struct {
	I
}

func (S1) show(s string) {
	println("show s1")
}

func (S1) say(a int) {
	println("say s1")
}

func (s1 *S1) SetAge(age int) {
	s1.age = age
}

// -------------------------------------------------------

func (c *S2) show(s string) {
	println("show *s2", s, c)
}

func (c S2) say(a int) {
	println("say s2", a, &c)
}

func TestInterfaceWithCommonTypeOrPointerType(t *testing.T) {
	var (
		i1 I = S1{}
		i2 I = &S2{i1}
		i3 I
	)
	i2.show("a")
	i2.say(0)
	ss := i2.(*S2)
	ss.show("b")
	ss.say(1)
	ss.I.show("c")
	ss.I.say(2)
	println("----------------------------------")

	println(i1 == i2, i3 == nil, (I)(nil))

	s2 := S2{}
	println(&s2)
	s2.say(5)
	s2.show("lu")

	sp := (*S2)(nil)
	fmt.Println(sp, &sp)

	var1 := I(S1{})
	var2 := I(&S2{})
	println(var1, var2)

	// NOTE 函数参数是个接口类型,直接传递其实现的实例也可以, 并不一定要传递 用接口引用的实例
	f1 := func(i I) { i.say(1994) }
	f2 := func() S1 {
		return S1{}
	}
	f1(f2())

}

// TODO struct 的接口实现和继承接口, 以及继承结构体
