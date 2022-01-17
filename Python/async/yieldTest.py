# -*- encoding: utf-8 -*-

import threading
from typing import Generator


def yieldCount(begin: int) -> int:
    for i in range(3):
        print(f"loop {i}")
        a = yield i  # 右边往上, 左边往下
        print(f"recv {a}")
        if a == "c":
            return begin
    return a


def testYieldCount():
    yc: Generator = yieldCount(0)
    for i in yc:
        print(i)

    # NOTE 生成器中的 return 只能起到中断生成器的作用
    # next(yc)
    print("----------------------------------")
    yc: Generator = yieldCount(0)
    print("return", yc.send(None))  # 首次不会走左边的赋值语句,所以规定设置一个None
    print("return", yc.send("a"))
    print("return", yc.send("b"))


def gen_func():
    a = yield 1
    print("a: ", a)
    b = yield 2
    print("b: ", b)
    c = yield 3
    print("c: ", c)
    return 4


def middle():
    gen = gen_func()
    ret = yield from gen  # 连接器 通道 ? 优化了上边的无法接收return 值 ???
    print("ret: ", ret)
    # return "middle Exception"


def main():
    mid = middle()
    for i in range(4):
        if i == 0:
            print(mid.send(None))
        else:
            try:
                print(mid.send(i))
            except StopIteration as e:
                print("e: ", e)


if __name__ == '__main__':
    print(threading.main_thread().getName())
    main()

# NOTE 以下讨论是基于linux 平台实现的 ?
# TODO python 中基于协程和IO多路复用的异步IO 到底是快在哪里 ?
# 事件循环驱动协程处理IO Coroutine
