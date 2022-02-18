package abase

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

type Human struct {
	name  string
	age   int
	phone string
}
type Student struct {
	Human  //匿名字段
	school string
	loan   float32
}
type Employee struct {
	Human   //匿名字段
	company string
	money   float32
}

func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

func (e Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name, e.company, e.phone)
}

type Men interface {
	SayHi()
	Sing(lyrics string)
}

type Woman interface {
	// NOTE 接口的 嵌入,注意不能与 java或者pytohn 的类继承混为一谈
	Men
	cry()
}

// -------------------------------------------------------------------------------------------------
type Controller struct {
	M int32
}

type Something interface {
	Get()
	Post()
}

func (c *Controller) Get() {
	fmt.Print("GET")
}

func (c *Controller) Post() {
	fmt.Print("POST")
}

type T struct {
	Controller
}

func (t *T) Get() {
	//new(test.Controller).Get()
	fmt.Print("T")
}
func (t *T) Post() {
	fmt.Print("T")
}

// -------------------------------------------------------------------------------------------------
// NOTE Go语言中的接口
func TestInterface(t *testing.T) {
	/*
		面向对象世界中的接口的一般定义是“接口定义对象的行为”。它表示让指定对象应该做什么。
		实现这种行为的方法(实现细节)是针对对象的.

		在Go中，接口是一组方法签名。当类型为接口中的 "NOTE 所有方法" 提供定义时，它被称为实现接口。
		它与OOP非常相似。接口指定了类型应该具有的方法，类型决定了如何实现这些方法.
		它把所有的具有共性的方法定义在一起，任何其他类型只要实现了这些方法就是实现了这个接口.
		接口定义了一组方法，如果某个对象实现了某个接口的所有方法，则此对象就实现了该接口。
	*/

	// interface可以被任意的对象实现
	// 一个对象可以实现任意多个interface
	// NOTE 任意的类型都实现了空interface(定义：interface{})，也就是包含0个method的interface.

	paul := Human{"Paul", 26, "111-222-XXX"}
	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	Tom := Employee{Human{"Sam", 36, "444-222-XXX"}, "Things Ltd.", 5000}

	// NOTE 多态
	//      interface的变量可以引用实现这个interface的任意类型的实例对象, 但是该引用变量,不能调用实现对象的属性.
	// ??? 使用指针的方式，也是可以的 ????

	var i Men
	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")
	i = Tom
	fmt.Println("This is Tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	x[0], x[1], x[2] = paul, Tom, mike
	for _, value := range x {
		value.SayHi()
	}

	// ---------------------------------------------------------------------------------------------
	/*
		NOTE interface函数参数
		Controller实现了所有的Something接口方法，当结构体T中调用Controller结构体的时候，T就相当于Java中的继承.
		T继承了Controller，因此，T可以不用重写所有的Something接口中的方法，因为父构造器已经实现了接口。
		如果Controller没有实现Something接口方法，则T要调用Something中方法，就要实现其所有方法。
		如果something = new(Controller)则调用的是Controller中的Get方法。
		T可以使用Controller结构体中定义的变量

		当一个接口类型作为函数参数时,意味着你可以提供多种实现了接口的类型.
	*/
	println("---------------------------------------------------------------------------")
	var something Something = new(T)
	var tt T
	tt.M = 1
	// tt.Controller.M = 1
	something.Get()
}

/*
	If it looks like a duck, swims like a duck, and quacks like a duck, then it probably is a duck.
	Duck Typing，鸭子类型，是动态编程语言的一种对象推断策略，它更关注对象能如何被使用，而不是对象的类型本身。
	Go 语言作为一门静态语言，它通过接口的方式完美支持鸭子类型。
	而在静态语言如 Java, C++ 中，必须要显示地声明实现了某个接口，之后，才能用在任何需要这个接口的地方。

	如果你在程序中调用某个数，却传入了一个根本就没有实现另一个的类型，那在编译阶段就不会通过。
	这也是静态语言比动态语言更安全的原因。

	静态语言在编译期间就能发现类型不匹配的错误，不像动态语言，必须要运行到那一行代码才会报错。
	当然，静态语言要求程序员在编码阶段就要按照规定来编写程序，为每个变量规定数据类型，
	这在某种程度上，加大了工作量，也加长了代码量。
	动态语言则没有这些要求，可以让人更专注在业务上，代码也更短，写起来更快，这一点，写 python 的同学比较清楚。

	Go 语言作为一门现代静态语言，是有后发优势的。它引入了动态语言的便利，同时又会进行静态语言的类型检查。
	NOTE Go 采用了折中的做法：不要求类型显示地声明实现了某个接口，只要实现了相关的方法即可，编译器就能检测到。

	鸭子类型是一种动态语言的风格，一个对象有效的语义，不是由继承自特定的类或实现特定的接口来决定，
	而是由它"当前方法和属性的集合"决定。
	NOTE Go 作为一种静态语言，通过接口实现了鸭子类型，实际上是 Go 的编译器在其中作了隐匿的转换工作。

	Go语言的多态性：
		Go中的多态性是在接口的帮助下实现的，接口可以在Go中隐式地实现。
		如果类型为接口中声明的所有方法提供了定义，则实现一个接口。
		任何定义接口所有方法的类型都被称为隐式地实现该接口。
		类型接口的变量可以保存实现接口的任何值。接口的这个属性用于实现Go中的多态性。
*/

// NOTE Go语言中接口类型的断言
func TestInterfaceAssert(t *testing.T) {
	// var i1 interface{} = new(Student)
	var i1 interface{} = Student{}
	s := i1.(Student)     //不安全，如果断言失败，会直接panic
	s, ok := i1.(Student) //安全，断言失败，也不会panic，只是ok的值为false
	if ok {
		fmt.Println(s)
	}

	// cannot type switch on non-interface value s (type Student)go
	// switch ins := s.(type) {
	switch ins := i1.(type) {
	case Human:
		fmt.Println("人", ins.name)
	case Student:
		fmt.Println("学生", ins.age)
	case Employee:
		fmt.Println("雇员", ins.company)
	default:
		fmt.Println("不是人")
	}

	type p struct{}

	var pp interface{} = p{}

	_, ok2 := pp.(*int) // 类型转换为一个整形的指针类型
	print(ok2)

	// ok1 := pp.(int) // 类型转换为一个整形
	// print(ok1) 不接收第二个变量,如果转换失败,则直接panic
}

// NOTE Go语言中接口 关键字 type
func TestType(t *testing.T) {
	// NOTE 定义 结构体
	// NOTE 定义 接口
	// NOTE 定义 新类型
	type myint int
	i := 1
	var mi myint = 2
	// i = mi
	// NOTE cannot use mi (type myint) as type int in assignment
	fmt.Println(i, mi)

	// NOTE 定义 函数类型 如果函数的类型比较复杂，使用type提前定义这个函数的类型
	type my_fun func(int, int) string

	f := func() my_fun {
		fun := func(a, b int) string {
			s := strconv.Itoa(a) + strconv.Itoa(b)
			return s
		}
		return fun
	}
	fmt.Println(f()(1, 2))

	// NOTE 定义类型别名 TypeAlias 只是 Type 的别名，本质上 TypeAlias 与 Type 是同一个类型
	// 在 C/C++语言中，代码重构升级可以使用宏快速定义新的一段代码,但是Go没有引入 宏,而是通过类型别名来解决.

	// Go 1.9 版本之前的内建类型定义 被改为新建类型
	// type byte uint8
	// Go 1.9 版本之后,并不是新建类型,而是起别名,以前版本的起别名,改为新建类型
	type ccc = uint8

	// NOTE 非本地类型不能定义方法
	// 定义time.Duration的别名为MyDuration
	type MyDuration = time.Duration
	// 为MyDuration添加一个函数
	EasySet := func(m MyDuration) (a string) {
		//cannot define new methods on non-local type time.Duration
		return ""
	}
	EasySet(12)
	// 不能在一个非本地的类型 time.Duration 上定义新方法。非本地方法指的就是使用 time.Duration 的代码所在的包.
	// 因为 time.Duration 是在 time 包中定义的, time.Duration 包与 当前源码包不在同一个包中，
	// 因此不能为不在一个包中的类型定义方法。

	// 将类型别名改为类型定义： type MyDuration time.Duration，也就是将 MyDuration 从别名改为类型。
	// 将 MyDuration 的别名定义放在 time 包中。

	// NOTE 类型别名作为结构体嵌入的成员时,如果别名指向同一类型,则会出现 ambiguous 错误
}
