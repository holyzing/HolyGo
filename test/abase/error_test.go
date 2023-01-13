package abase

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

func TestError(t *testing.T) {
	/*
		Go语言没有提供像Java、C#语言中的try...catch异常处理方式，而是通过函数返回值逐层往上抛。
		这种设计，鼓励工程师在代码中显式的检查错误，而非忽略错误，好处就是避免漏掉本应处理的错误。
		但是带来一个弊端，让代码啰嗦

		Java 存在运行时和编译时异常,作为编译型语言,GO语言是无法捕获编译性异常的,只能直接修复该异常.

		Java中规定异常或者错误必须是 Throwable（子类 Error和Exception）的，且只能 throws 或者 try...catch...，且不能作为值赋值给变量。
		对于未抛出或者未捕获的异常，异常发生后，如果找不到处理器，就终止执行程序，在控制台打印出 stack trace。

		而Go只是在约定中将引用类型interface Error类型作为异常标志，且是作为值在函数中传递的，它就是一个一般意义上的类型的值而已。
		处理错误的惯用方法是将返回的错误与nil进行比较。nil值表示没有发生错误，而非nil值表示出现错误。

		实现该接口方法的类型 都属于 error 类型, 当打印该类型时, fmt.Println 会调用 Error 方法获取错误描述.
		但是 单单通过用户态的描述信息,有时候很难定位错误,更需要相关的堆栈信息来定位,分析错误.
			type error interface {
		    	Error() string
			}
	*/

	// NOTE 1.断言底层结构类型并从结构字段获取更多信息
	f, err := os.Open("/test.txt")
	if err, ok := err.(*os.PathError); err == nil {
		fmt.Println(f.Name(), "opened successfully")
	} else if ok {
		fmt.Println("File at path", err.Path, "failed to open")
	} else {
		fmt.Println("Uncaught error")
	}

	// NOTE 2. 断言底层结构类型，并通过调用struct类型的方法获取更多信息
	addr, err := net.LookupHost("golangbot123.com")
	if err, ok := err.(*net.DNSError); err == nil {
		fmt.Println(addr)
	} else if ok {
		if err.Timeout() {
			fmt.Println("operation timed out")
		} else if err.Temporary() {
			fmt.Println("temporary error")
		} else {
			fmt.Println("generic error: ", err)
		}
		return
	}

	// NOTE 3. 直接与类型错误的变量进行比较
	// error 类型 定义了 == 操作符
	// Glob函数用于返回与模式匹配的所有文件的名称。
	// 当模式出现错误时，该函数将返回一个错误ErrBadPattern
	files, err := filepath.Glob("[")
	if err == filepath.ErrBadPattern {
		fmt.Println(err)
	} else {
		fmt.Println("matched files", files)
	}

	// NOTE 4. 不要号忽略错误
	files2, _ := filepath.Glob("[")
	fmt.Println("matched files", files2)
	// 使用行号中的空白标识符,这个模式本身是畸形的, 由于忽略了这个错误，输出看起来好像没有文件匹配这个模式
	// 这会造成一个错误的业务判断

	// NOTE 5. GO语言中的自定义错误
	/*
		NOTE 使用struct类型和字段提供关于错误的更多信息
		还可以使用将错误接口实现为错误的struct类型。这给我们提供了更多的错误处理的灵活性。
		如果想要访问导致错误的具体 “数据”，如果单纯的从 error.Error() 中解析错误描述，
		这不是一种正确的方法，因为如果描述发生了变化，那么我们解析字符串的代码就会中断。
		模仿标准库异常信息的策略，在“断言底层结构类型并从struct字段获取更多信息”，
		并使用struct字段来提供对导致错误的“数据”的访问。
		创建一个实现错误接口的struct类型，并使用它的字段来提供关于错误的更多信息。
	*/
}

// -------------------------------------------------------------------------------------------------

/*
	Golang中引入两个内置函数 panic 和 recover 来触发和终止异常处理流程，
	同时引入关键字defer来延迟执行defer后面的函数。

	一直等到包含defer语句的函数执行完毕时，延迟函数（defer后的函数）才会被执行，
	而不管包含defer语句的函数是通过return的正常结束，还是由于panic导致的异常结束。
	可以在一个函数中执行多条defer语句，它们的执行顺序与声明顺序相反。

	当程序运行时，如果遇到引用空指针、下标越界或显式调用panic函数等情况，
	则先触发panic函数的执行，然后调用延迟函数。调用者继续传递panic，因此该过程一直在调用栈中重复发生：
	函数停止执行，调用延迟执行函数等。如果一路在延迟函数中没有recover函数的调用，则会到达该协程的起点，
	该协程结束，然后终止其他所有协程，包括主协程（类似于C语言中的主线程，该协程ID为1）。

NOTE panic
	1、内建函数
	2、假如函数F中书写了panic语句，会终止其后要执行的代码，
	   在panic所在函数F内如果存在要执行的defer函数列表，按照defer的逆序执行
	3、返回函数F的调用者G，在G中，调用函数F语句之后的代码不会执行，
	   假如函数G中存在要执行的defer函数列表，按照defer的逆序执行，
	   这里的defer 有点类似 try-catch-finally 中的 finally
	4、同3依次退出函数调用链，直到goroutine整个退出，并报告错误

NOTE recover
	1、内建函数
	2、用来控制一个goroutine的panicking行为，捕获panic，从而影响应用的行为
	3、一般的调用建议
		a). 在defer函数中，通过recever来终止一个goroutine的panicking过程，从而恢复正常代码的执行
		b). 可以获取通过panic传递的error

NOTE 简单来讲：go中可以抛出一个panic的异常，然后在defer中通过recover捕获这个异常，然后正常处理。

NOTE 错误和异常从Golang机制上讲，就是error和panic的区别。
     在C++/Java，没有error但有errno，没有panic但有throw。

NOTE Golang错误和异常是可以互相转换的：
	 错误转异常，比如程序逻辑上尝试请求某个URL，最多尝试三次，
	 尝试三次的过程中请求失败是错误，尝试完第三次还不成功的话，失败就被提升为异常了。

	 异常转错误，比如panic触发的异常被recover恢复后，
	 将返回值中error类型的变量进行赋值，以便上层函数继续走错误处理流程。

NOTE 什么情况下用错误表达，什么情况下用异常表达，就得有一套规则，否则很容易出现一切皆错误或一切皆异常的情况。

	异常处理的作用域（场景）：
		空指针引用
		下标越界
		除数为0
		不应该出现的分支，比如default
		输入不应该引起函数错误

	其他场景使用错误处理，这使得我们的函数接口很精炼。
	对于异常，我们可以选择在一个合适的上游去recover，并打印堆栈信息，使得部署后的程序不会终止。

NOTE Golang错误处理方式一直是很多人诟病的地方，有些人吐槽说一半的代码都是
     "if err != nil { / 打印 && 错误处理 / }"，严重影响正常的处理逻辑。
	 当我们区分错误和异常，根据规则设计函数，就会大大提高可读性和可维护性。
*/

// NOTE 错误处理的正确姿势
// 姿势一：失败的原因只有一个时，不使用error
// 姿势二：没有失败时，不使用error
// 姿势三：error应放在返回值类型列表的最后,对于返回值类型error，用来传递错误信息，通常放在最后一个。
// 姿势四：错误值统一定义，而不是跟着感觉走,在Golang的每个包中增加一个错误对象定义文件.保证同一错误信息复用。
// 姿势五：错误逐层传递时，层层都加日志
// 姿势六：错误处理使用defer
func DeferErrorTest() {
	// NOTE 当Golang的代码执行时，如果遇到defer的闭包调用，则压入堆栈。
	//      当函数返回时，会按照后进先出的顺序调用闭包。
	//      对于闭包的参数是值传递，而对于外部变量却是引用传递，
	//      NOTE 所以闭包中的外部变量err的值就变成外部函数返回时最新的err值
	/*
		err := createResource1()
		if err != nil {
			return "ERR_CREATE_RESOURCE1_FAILED"
		}
		defer func() {
			if err != nil {
				// 引用传递，defer函数执行的时候，引用的是最后一次赋值的 error，如果存在参数，则是值传递
				destroyResource1()
			}
		}()
		err = createResource2()
		if err != nil {
			return "ERR_CREATE_RESOURCE2_FAILED"
		}
		defer func() {
			if err != nil {
				destroyResource2()
			}
		}()
		err = createResource3()
		// ...
	*/
}

// 姿势七：当尝试几次可以避免失败时，不要立即返回错误,
//        比如网络请求等依赖外部条件的操作，但要设置重试次数以及重试间隔时间

// 姿势八：当上层函数不关心错误时，建议不返回error
// 		  对于一些资源清理相关的函数（destroy/delete/clear），如果子函数出错，打印日志即可，
// 		  而无需将错误进一步反馈到上层函数，因为一般情况下，上层函数是不关心执行结果的，
// 		  或者即使关心也无能为力，建议将相关函数设计为不返回error。

// 姿势九：当发生错误时，不忽略有用的返回值
//        通常，当函数返回non-nil的error时，其他的返回值是未定义的(undefined)，
//        这些未定义的返回值应该被忽略。然而，有少部分函数在发生错误时，仍然会返回一些有用的返回值。
//        比如，当读取文件发生错误时，Read函数会返回可以读取的字节数以及错误信息。
//        对于这种情况，应该将读取到的字符串和错误信息一起打印出来。
//        对函数的返回值要有清晰的说明，以便于其他人使用。

// NOTE 异常处理的正确姿势
// 姿势一：在程序开发阶段，坚持速错
//        速错，简单来讲就是“让它挂”，只有挂了你才会第一时间知道错误。在早期开发以及任何发布阶段之前，
//        最简单的同时也可能是最好的方法是调用panic函数来中断程序的执行以强制发生错误，
//        使得该错误不会被忽略，因而能够被尽快修复。

// 姿势二：在程序部署后，应恢复异常避免程序终止
//        在Golang中，某个Goroutine如果panic了，并且没有recover，那么整个Golang进程就会异常退出。
//        所以，一旦Golang程序部署后，在任何情况下发生的异常都不应该导致程序异常退出，
//        我们在上层函数中加一个延迟执行的recover调用来达到这个目的，
//        并且是否进行recover需要根据环境变量或配置文件来定，默认需要recover。
//        这个姿势类似于C语言中的断言，但还是有区别：一般在Release版本中，断言被定义为空而失效，
//        但需要有if校验存在进行异常保护，尽管契约式设计中不建议这样做。
//        在Golang中，recover完全可以终止异常展开过程，省时省力。

//   	  调用recover的延迟函数中以最合理的方式响应该异常：
//		  打印堆栈的异常调用信息和关键的业务信息，以便这些问题保留可见；
//		  将异常转换为错误，以便调用者让程序恢复到健康状态并继续安全运行。

func TestPanicRecover(t *testing.T) {
	println("--------------------------------------------")
	funcB := func() string {
		println("funcB")
		if 2 > 1 {
			panic("foo")
		}
		println("funcB before return")
		return "dasd"
		// return errors.New("success")  // unreachable code
	}

	funcA := func(err error) string {
		println("funcA")
		defer func() {
			println("funcA defer")
			if p := recover(); p != nil {
				fmt.Printf("panic recover! p: %v\n", p)
				str, ok := p.(string)
				if ok {
					err = errors.New(str)
				} else {
					err = errors.New("panic")
				}
				debug.PrintStack()
			}
		}()
		return funcB() // return 是非原子的
	}

	// 原因是panic异常处理机制不会自动将错误信息传递给error，所以要在funcA函数中进行显式的传递
	var err error
	res := funcA(err)
	// ??? 如果在一个函数内遇到panic，在函数外部recover，则该函数最后的返回值会是它的零值
	if res != "" {
		fmt.Printf("err is nil\n")
	} else {
		fmt.Printf("err is %v\n", err)
	}
	println("--------------------------------------------")
}

// 姿势三：对于不应该出现的分支，使用异常处理

// 姿势四：针对入参不应该有问题的函数，使用panic设计
// 	      入参不应该有问题一般指的是硬编码，例：MustCompile函数是对Compile函数的包装。

// func MustCompile(str string) *Regexp {
//     regexp, error := Compile(str)
//     if error != nil {
//         panic(regexp: Compile( + quote(str) + ):  + error.Error())
//     }
//     return regexp
// }

// 所以，对于同时支持用户输入场景和硬编码场景的情况，一般支持硬编码场景的函数是对支持用户输入场景函数的包装。
// 对于只支持硬编码单一场景的情况，函数设计时直接使用panic，即返回值类型列表中不会有error，
// 这使得函数的调用处理非常方便（没有了乏味的"if err != nil {/ 打印 && 错误处理 /}"代码块）。
