package main

import "fmt"

// private， used in under the package
func appInit() {
	fmt.Println("--- This is a initial app ！ ---")
}

// public
func AppInit() {
	appInit()
}

func init() {
	println("同一个包下 按照源码文件名称的字符顺序")
}
