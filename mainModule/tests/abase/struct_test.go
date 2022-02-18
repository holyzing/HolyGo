package abase

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	// "./advanced"
	"mainModule/tests/advanced"
)

type People struct {
	name string
	age  int
}

func (p People) show() string {
	// NOTE string()会直接把字节或者数字转换为字符的UTF-8表现形式
	println("---------- instance show ")
	mesg := "name:" + p.name + " " + "age:" + strconv.Itoa(p.age)
	print(mesg)
	return mesg
}

// ??? 结构体指针变量和结构体变量为什么作为函数或者方法参数,在函数或者方法内部都可以 . (点)属性
func (p *People) showPointer() string {
	// NOTE string()会直接把字节或者数字转换为字符的UTF-8表现形式
	println("---------- pointer show ")
	mesg := "name:" + p.name + " " + "age:" + strconv.Itoa(p.age) + (*p).name
	print(mesg)
	return mesg
}

func (p *People) calledByUninstantiateiStruct() {
	println("called by struct with a nil struct")
}

// NOTE 引用传递，值传递，引用类型，值类型
func TestCalledByUninstantiateiStruct(t *testing.T) {
	var p People // 结构体本身的一个拷贝 ????
	// var p3 = People
	p.age = 20
	p.calledByUninstantiateiStruct()
	p.show()
	p.showPointer()

	var i int     // default value： 0
	var s string  // default value: ""
	var p3 People // default value: People{}
	println(p.age, p3.age, &p.age, &p3.age, i, s)

	println("###################################################")

	var p2 *People // default value: nil
	p2.calledByUninstantiateiStruct()
	p2 = &p
	p2.show() // 接收一个类型实例的方法不能传入一个空指向的指针
	// panic: runtime error: invalid memory address or nil pointer dereference
	p2.showPointer()

	// 指针作为函数参数，其实也是值传递，但是指针的操作符.被重载为访问指针指向的内存地址
	// 但是go 应该是分引用类型和值类型的，这两种类型的定义是？还是人们的约定俗称 ？
	println("###################################################")
	println(&p2, p2.age)
	func(p *People) {
		println(&p, p.age)
	}(p2)
}

// NOTE GO语言中的结构体
func TestStruct(t *testing.T) {
	type LocalPerson struct {
		name string
		age  int
	}

	// NOTE 局部(函数)中不能声明方法
	// func (p Persion) show(){
	// }
	// func test2(){
	// }

	p1 := LocalPerson{name: "ll", age: 22}
	p2 := LocalPerson{}
	p3 := new(LocalPerson)
	fmt.Println(p1, p2, p3)

	// ??? 结构体和数组(切片)的成员(元素访问) 均可通过其变量或者指针变量直接访问 ??? 这是为什么呢

	pp1 := &p1
	fmt.Println(p1.name, pp1.name)

	arr := [3]int{}
	arr[0] = 1
	arrp := &arr
	arrp[1] = 2
	fmt.Println(arr)

	pp := People{}
	fmt.Printf("%v\n", pp)
	pp.show()

	type Student struct {
		People // NOTE 继承匿名字段类型的属性
		name   string
		grade  string
	}
	println("--------------------------------------------------------------------------------")
	// mixture of field:value and value initializers
	// s := Student{People{name: "zuz", age: 20}, name: "sub", grade: "ernian"}
	s := Student{People{name: "zuz", age: 20}, "sub", "ernian"}
	fmt.Printf("%v\n", s)
	fmt.Println(s.name, s.People.name, s.show(), s.showPointer())
	s.People = People{}
	s.age = -1
	fmt.Printf("%v\n", s.People)

	// NOTE 不仅仅是struct字段哦，所有的内置类型和自定义类型都是可以作为匿名字段的
	f := First{name: 78}
	fmt.Println(f.name)
	// NOTE 同一包下访问私有
	fmt.Printf("%v\n", second{name: 23, a: 1, b: false})

	type Dog struct {
		Student Student
	}
	var dog Dog
	dog.Student = Student{s.People, "sub2", "yinian"}
	fmt.Printf("%v\n", dog)

	dog2 := Dog{Student{People{}, "sub3", "chuoxue"}}
	fmt.Printf("%v\n", dog2)
	println("--------------------------------------------------------------------------------")

	// NOTE 在结构体中属于匿名结构体的字段称为提升字段，因为它们可以被访问，
	//      就好像它们属于拥有匿名结构字段的结构一样。该种现象称为字段提升?
	//      就比如 Student 可以直接访问 People 的字段.

	// NOTE 如果结构体类型以大写字母开头，那么它是一个导出类型，可以从其他包访问它。
	//      类似地，如果结构体的字段以大写开头，则可以从其他包访问它们。
	//      比如在 more 包中使用 父包中定义 的 Second

	// cannot refer to unexported field 'b' in struct literal of type advanced.Morego
	m := advanced.More{A: 97}
	fmt.Println(m.A) // 只能访问到 A
}

// NOTE Go语言中的结构体
func TestMoreStruct(t *testing.T) {
	// NOTE 结构体是值类型，如果每个字段具有可比性，则是可比较的。
	//      如果它们对应的字段相等，则认为两个结构体变量是相等的。
	type name struct {
		firstName string
		lastName  string
	}
	name1 := name{"Steve", "Jobs"}
	name2 := name{"Steve", "Jobs"}
	if name1 == name2 {
		fmt.Println("name1 and name2 are equal")
	} else {
		fmt.Println("name1 and name2 are not equal")
	}

	// type image struct {
	// 	data map[int]int
	// }
	// image1 := image{data: map[int]int{
	// 	0: 155,
	// }}
	// image2 := image{data: map[int]int{
	// 	0: 155,
	// }}
	// invalid operation: image1 == image2 (struct containing map[int]int cannot be compared)
	// if image1 == image2 {
	// 	fmt.Println("image1 and image2 are equal")
	// }
	// fmt.Println(image1, image2)
	println("--------------------------------------------")

	show := func(x name) {
		x.firstName = "cao"
		// ??? 局部编译后,并加载到内存后,或者形成缓存,x 的地址在局部内(函数体内),无论调用几次都不变????
		println(&x)    // TODO 注意该函数参数的格式要求
		fmt.Println(x) // TODO 注意它与 println 的区别
	}
	im := name{firstName: "fi1", lastName: "la1"} // NOTE 直接构建 是一个 普通变量 name
	imp := new(name)                              // NOTE new 返回的是一个指针变量 *name
	imp.firstName = "fi2"
	imp.lastName = "la2"

	fmt.Println(reflect.TypeOf(im), reflect.TypeOf(imp))
	println(&im, imp)
	show(im)   // 值拷贝
	show(*imp) // 值拷贝
	fmt.Println(im, *imp)

	// ??? 只有指针作为参数时才是 引用传递 ??? 取值符 & 取的是变量所指向内存的地址 ?? 还是变量的内存, 地址 ???
	println("--------------------------------------------")
}

/*
new用于各种类型的内存分配, new 本质上跟其它语言中的 new 功能基本一样：
new(T)分配了零值填充的T类型的内存空间，并且返回其地址，即一个*T类型的值。
用Go的术语说，它返回了一个指针，指向新分配的类型T的零值。

make用于内建类型（map、slice 和channel）的内存分配。
make(T, args)与new(T)有着不同的功能，make只能创建slice、map和channel，
并且返回一个有初始值(非零)的T类型，而不是*T。
本质来讲，导致这三个类型有所不同的原因是指向数据结构的引用在使用前必须被初始化。
例如，一个slice，是一个包含指向数据（内部array）的指针、长度和容量的三项描述符；
在这些项目被初始化之前，slice为nil。对于slice、map和channel来说，make初始化了内部的数据结构，填充适当的值。
*/
