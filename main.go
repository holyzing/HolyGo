package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	// _ 操作其实是引入该包，而不直接使用包里面的函数，而是调用了该包里面的init函数。
	//   也就是说，使用下划线作为包的别名，会仅仅执行init().
)

// 对同一个 package 中的不同文件，将文件名按字符串进行“从小到大”排序，之后顺序调用各文件中的init()函数。
// 对于不同的 package，如果不相互依赖的话，按照 main 包中 import 的顺序调用其包中的 init() 函数。
// 如果 package 存在依赖，调用顺序为最后被依赖的最先被初始化, main 包总是被最后一个初始化，它总是依赖别的包。
// 一个包被其它多个包 import，但只能被初始化一次。

// 也就是说，对于go语言来讲，其实并不关心你的代码是内部还是外部的，总之都在GOPATH里，任何import包的路径都
// 是从GOPATH开始的；唯一的区别，就是内部依赖的包是开发者自己写的，外部依赖的包是go get下来的。

// 一个非main包在编译后会生成一个.a文件（在临时目录下生成，除非使用go install安装到$GOROOT或$GOPATH下，
// 否则看不到.a），用于后续可执行程序链接使用。

// 例如Go标准库中的包对应的源码部分路径在：$GOROOT/src，而标准库中包编译后的
// .a文件路径在$GOROOT/pkg/darwin_amd64下

func init() {
	fmt.Println("-------- main: init1 --------")
}

func init() {
	fmt.Println("-------- main: init2 --------")
}

// func main and int must have no arguments and no return values
func main() {

	fmt.Println("-------- main: main --------")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// collects.Cao()
	err := r.Run()
	if err != nil {
		return
	}
}

// go env -w GOPROXY=https://goproxy.cn,direct
// 在当前目录下检测到 go.mod 就会以 Module111 模式检测依赖

// Module： 顶层文件夹init， 子文件夹不应该作为一个 Module ？

// 1. 同一模块 同一包下 的所有文件内直接定义的函数，结构体等“元素” 可直接互相引用，
//    不需要导入，也不需要区分大小写，但是注意不要造成循环引用。
// 2. 一个包下的所有源文件只能引用同一个包名，但不一定要与包文件夹名称相同
// 3. 当前包下执行 main 函数时，引用本地其它源文件 需要全部参与编译 即 *.go
// 4. 绝对导入 以模块开头引用开始

// 5. GO 中只能导入包，然后直接引用包中定义的元素
// 6. import( . “fmt” )
// 7. "database/sql" _
// 8. go run *.go 不能运行 _test.go

// 库源码文件被安装后，相应的归档文件（.a 文件）会被存放到当前工作区的 pkg 的平台相关目录下。
// 名称以 _test.go 为后缀的代码文件，并且必须包含 Test 或者 Benchmark 名称前缀的函数
// 名称以 Test 为名称前缀的函数，只能接受 *testing.T 的参数，这种测试函数是功能测试函数。
// 名称以 Benchmark 为名称前缀的函数，只能接受 *testing.B 的参数，这种测试函数是性能测试函数。

/**
NOTE vender
*/
