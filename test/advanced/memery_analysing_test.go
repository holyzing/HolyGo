package advanced

import (
	"fmt"
	"runtime"
	"testing"
)

// NOTE 逃逸
func TestMemoryEscape(t *testing.T) {
	fmt.Println(runtime.NumCPU())
	i := 1
	i++
	i--
}

// 函数内部的变量一般是分配在栈中的，函数返回的局部变量是局部变量的拷贝,函数结束后，栈空间被收回
// 当函数内部返回一个局部变量的地址，则当函数结束后，局部变量被回收，在函数外部访问变量的地址，会出现段错误，
// 此时局部变量需要定义在堆中。

// 但是在go语言中，函数是应该定义在堆上还是栈中，由编译器来进行逃逸分析后决定。


// NOTE 注意：go 在编译阶段确立逃逸，并不是在运行时。

// WHAT 函数在编译的时候就能确定该分配多大的栈内存 ???

// “生命期会超过当前stack frame的指针都会逃逸”，比如存到全局变量，通过接口或函数变量调用函数时的指针参数(
// 因为编译期无法它们到底绑定了哪个函数或方法, 也就无法确定会它们会怎么使用参数)。分析问题时还是得靠编译器输
// 出逃逸分析结果，不是靠程序员背八股瞎猜，但是掌握了基本原则写代码时就能避免无谓的逃逸，
// 比如可能很多人不知道fmt.Print系列就都会逃逸……

// 内存逃逸与 闭包

/**
	栈 可以简单得理解成一次函数调用内部申请到的内存，它们会随着函数的返回把内存还给系统。
	函数内部申请的临时变量，即使你是用make申请到的内存，如果发现在退出函数后没有用了，那么就把丢到栈上，毕竟栈上的内存分配比堆上快很多。

	函数内部变量申请后作为返回值返回了，编译器会认为在退出函数之后还有其他地方在引用，当函数返回之后并不会将其内存归还。那么就申请到堆里。

	如果变量都分配到堆上，堆不像栈可以自动清理。它会引起Go频繁地进行垃圾回收，而垃圾回收会占用比较大的系统开销。

	堆适合不可预知的大小的内存分配。但是为此付出的代价是分配速度较慢，而且会形成内存碎片。

	栈内存分配则会非常快，栈分配内存只需要两个CPU指令：“PUSH”和“RELEASE”分配和释放；而堆分配内存首先需要去找到一块大小合适的内存块。
	之后要通过垃圾回收才能释放。

	逃逸分析是一种确定指针动态范围的方法。简单来说就是分析在程序的哪些地方可以访问到该指针。
	简单来说，编译器会根据变量是否被外部引用来决定是否逃逸：
		1、如果函数外部没有引用，则优先放到栈中(stack frame 栈帧)；
		2、如果函数外部存在引用，则必定放到堆中；
	对此你可以理解为，逃逸分析是编译器用于决定变量分配到堆上还是栈上的一种行为。

	提问：函数传递指针真的比传值效率高吗？
		我们知道传递指针可以减少底层值的拷贝，可以提高效率，但是如果拷贝的数据量小，
		由于指针传递会产生逃逸，可能会使用堆，也可能会增加GC的负担，所以传递指针不一定是高效的。

	From a correctness standpoint, you don’t need to know.
	Each variable in Go exists as long as there are references to it.
	The storage location chosen by the implementation is irrelevant to the semantics of the language.

	The storage location does have an effect on writing efficient programs.
	When possible, the Go compilers will allocate variables that are local to a function in that function’s stack frame.

	However, if the compiler cannot prove that the variable is not referenced after the function returns,
	then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors.
	Also, if a local variable is very large, it might make more sense to store it on the heap rather than the stack.

	In the current compilers, if a variable has its address taken,that variable is a candidate for allocation on the heap.
	However, a basic escape analysis recognizes some cases when such variables will not live past the return from the function and can reside on the stack.
 */


// 指针逃逸: 返回变量是个指针
// 栈空间不足: 当切片长度扩大到10000时就会逃逸。实际上当栈空间不足以存放当前对象时或无法判断当前切片长度时会将对象分配到堆中。
// 动态类型逃逸: 很多函数参数为interface类型, 比如 fmt.* 系列的函数 都有参数 a ...interface{}
//			   编译期间很难确定其参数的具体类型，也能产生逃逸。

/**
逃逸分析的作用是什么呢？
	1、逃逸分析的好处是为了减少gc的压力，不逃逸的对象分配在栈上，当函数返回时就回收了资源，不需要gc标记清除。
	2、逃逸分析完后可以确定哪些变量可以分配在栈上，栈的分配比堆快，性能好(逃逸的局部变量会在堆上分配 ,而没有发生逃逸的则有编译器在栈上分配)。
	3、同步消除，如果你定义的对象的方法上有同步锁，但在运行时，却只有一个线程在访问，此时逃逸分析后的机器码，会去掉同步锁运行。

总结
	1、堆上动态分配内存比栈上静态分配内存，开销大很多。
	2、变量分配在栈上需要能在编译期确定它的作用域，否则会分配到堆上。
	3、Go编译器会在编译期对考察变量的作用域，并作一系列检查，如果它的作用域在运行期间对编译器一直是可知的，那么就会分配到栈上。
       简单来说，编译器会根据变量是否被外部引用来决定是否逃逸。
	4、对于Go程序员来说，编译器的这些逃逸分析规则不需要掌握，我们只需通过go build -gcflags '-m’命令来观察变量逃逸情况就行了。
	5、不要盲目使用变量的指针作为函数参数，虽然它会减少复制操作。
       但其实当参数为变量自身的时候，复制是在栈上完成的操作，开销远比变量逃逸后动态地在堆上分配内存少的多。

		NOTE 比如函数的嵌套调用,就不要返回函数内部局部变量的指针了,直接返回该内部变量的拷贝就行了

	6、逃逸分析在编译阶段完成的。

	Go 语言的逃逸分析遵循以下两个不变性
		1. 指向栈对象的指针不能存在于堆中
		2. 指向栈对象的指针不能在栈对象回收后存活
 */



type Student struct {
	Name string
	Age  int
}

// go build -gcflags=-m
// NOTE 返回签名声明的变量,等价函数内部第一行的变量声明
func StudentRegister(name string, age int) (s *Student){
	s = new(Student) //局部变量s逃逸到堆
	s.Name = name
	s.Age = age
	return 
}

func TestStackEscape(t *testing.T) {
	StudentRegister("lu", 28)
}
