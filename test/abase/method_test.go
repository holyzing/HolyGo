package abase

import (
	"fmt"
	"math"
	"testing"
)

type Rectangle struct {
	width, height float64
}
type Circle struct {
	radius float64
}

func (r Rectangle) area() float64 {
	return r.width * r.height
}

func (r *Rectangle) setVal() {
	r.height = 20
}

func (c Circle) area() float64 {
	return c.radius * c.radius * math.Pi
}

func area(r Rectangle) float64 {
	return r.width * r.height
}

// NOTE area redeclared in this block previous declaration a
// func area(c Circle) float64 {
// 	return c.radius * c.radius * math.Pi
// }

// NOTE Go 语言中的 "成员方法"
func TestMethod(t *testing.T) {
	/*
		Go 语言中同时有函数和方法。
		一个方法就是一个包含了接受者的函数，接受者可以是命名类型或者结构体类型的一个值或者是一个指针。
		所有给定类型的方法属于该类型的方法集

		方法只是一个函数，它带有一个特殊的类型接收器，它是在func关键字和方法名之间编写的。
		接收器可以是struct类型或非struct类型。接收方可以在方法体内访问。

		方法能给用户自定义的类型添加新的行为。它和函数的区别在于方法有一个接收者，
		给一个函数添加一个接收者，那么它就变成了方法。
		接收者可以是值接收者，也可以是指针接收者。

		在调用方法的时候，值类型既可以调用值接收者的方法，也可以调用指针接收者的方法；
		指针类型既可以调用指针接收者的方法，也可以调用值接收者的方法。
		也就是说，不管方法的接收者是什么类型，该类型的值和指针都可以调用，不必严格符合接收者的类型。
	*/
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{10}
	c2 := Circle{25}
	fmt.Println("Area of r1 is: ", r1.area())
	fmt.Println("Area of r2 is: ", r2.area())
	fmt.Println("Area of c1 is: ", c1.area())
	fmt.Println("Area of c2 is: ", c2.area())

	// NOTE Java 中的方法重载,在一个类中不同参数列表的同名方法是允许存在的.
	// NOTE Java 中的方法重写,子类重写父类的相同签名的方法,会引出一些多态的问题

	// ??? 为什么我们可以用函数来写相同的程序呢?有以下几个原因
	// Go不是一种纯粹面向对象的编程语言，它不支持类。因此，类型的方法是一种实现类似于类的行为的方法。
	// 相同名称的方法可以在不同的类型上定义，而具有相同名称的函数是不被允许在同一作用域中(同一包内)定义的。
	// 同名的函数和方法也是可以共存的.
	area(r1)

	/**
	作用域为已声明标识符所表示的 (常量、类型、变量、函数或包) 在源代码中的作用范围。

	Go 语言中变量可以在三个地方声明：
		在函数体内声明的变量称之为局部变量，它们的作用域只在函数体内，参数和返回值变量也是局部变量。
		在函数体外声明的变量称之为全局变量，首字母大写全局变量可以在整个包甚至外部包（被导出后）使用。
		函数定义中的变量称为形式参数，形式参数会作为函数的局部变量来使用

	指针作为接收者
		NOTE 若不是以类型指针作为接收者，实际只是获取了一个实例copy，而不能真正改变传递给接收者实例中的数据,
			 即使它是类型实例的 "方法". 接收器是一个指针类型则会传递一个实例指针,是一个原类型,则会传递实例.
	*/

	p := Rectangle{1, 2}
	s := p // TODO 值拷贝
	p.setVal()
	fmt.Println(p.height, s.height)

	// NOTE method和成员变量一样是可以继承的，如果匿名字段实现了一个method，
	//      那么包含这个匿名字段的struct也能调用该method
	//      方法重写,按照就近原则,进行寻址调用.

	mark := Student{Human{"Mark", 25, "222-222-YYYY"}, "MIT", 2.1}
	sam := Employee{Human{"Sam", 45, "111-888-XXXX"}, "Golang Inc", 12121}
	mark.SayHi()
	sam.SayHi()
}

func TestBlockCode(t *testing.T) {
	arr := [5]int{}
	{
		arr[0] = 2
		var i int = 3
		fmt.Println(i)
	}
	// 局部代码块限定局部变量的作用域
	// fmt.Println(i)
}

// ----------------------------------------------------------------------

type TestPointerMethod struct {
	Name string
}

func (t *TestPointerMethod) ChangePointerName(newName string) {
	t.Name = newName
}

func (t TestPointerMethod) ChangeNonPointerName(newName string) {
	t.Name = newName
}

func TestMethodPointer(t *testing.T) {
	tpm := TestPointerMethod{Name: "origin"}
	fmt.Println(tpm)
	tpm.ChangePointerName("name1")
	fmt.Println(tpm)
	tpmp := &tpm
	tpmp.ChangePointerName("name2")
	fmt.Println(tpm)
	tpm.ChangeNonPointerName("name3")
	fmt.Println(tpm)
	tpmp.ChangeNonPointerName("name4")
	fmt.Println(tpm)
	// 指针方法，即使使用非指针引用，this也是该非指针变量的指针
	// 如果是非指针方法，则this，是该变量的值的拷贝，即使该变量是个指针变量
}
