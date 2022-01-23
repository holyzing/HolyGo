# -*- encoding: utf-8 -*-
"""
    1）使用最常规的同步编程要实现异步并发效果并不理想，或者难度极高。
    2）由于GIL锁的存在，多线程的运行需要频繁的加锁解锁，切换线程，这极大地降低了并发性能；

    1.协程是在单线程里实现任务的切换的
    2.利用同步的方式去实现异步
    3.不再需要锁，提高了并发性能

    1、调用方：调用委派生成器的客户端（调用方）代码
    2、委托生成器：包含yield from表达式的生成器函数
    3、子生成器：yield from后面加的生成器函数

    用户态，进行上下文切换
    一个线程内，切换执行逻辑

    # 多路复用，密集型 IO，执行IO 不占用 CPU逻辑吗 ？？？占用的少 ？？？
"""
import sys
from typing import Any, Generator  # , Iterator


def yield0(threshold):
    i = 11
    while i < threshold:
        print("--->yield0 before")
        r = yield i          # NOTE 产出并让步
        i += 1
        print("yield0 send:", r)
    # NOTE 执行到最后都没有遇到 yield 就会 抛出异常 StopIteration

# g = yield0(15)


def yieldFrom(threshold):
    i = 21
    g = yield0(15)
    # print(type(g), g.send(None))
    while i < threshold:
        print("--->yieldFrom before")
        r = yield from g           # NOTE 会自动预激 并且捕获 StopIteration 异常，停止产出
        i += 1
        print("yieldFrom send:", i, r)

# TODO 为什么 字符串 列表，字典，（生成器）（都属于 iterable的），一次 yield from 可以产出所有元素
#      而 由函数转化来的 生成器并不能一次 yield from 出所有元素 ？？
#      yield from 语法使用的限制 是 （函数体内 ？ 循环体内 ？）


def temp():
    # 字符串
    astr = 'ABC'
    # 列表
    alist = [1, 2, 3]
    # 字典
    adict = {"name": "wangbm", "age": 18}
    # 生成器
    agen: Generator[int, Any, None] = (i for i in range(4, 8))

    # <generator object temp.<locals>.<genexpr> at 0x7fd1d08bb2d0>
    # <class 'generator'>

    def gen(*args):
        for item in args:
            yield from item

    new_list = gen(astr, alist, adict, agen)
    print(list(new_list))


# 子生成器
def average_gen():
    total = 0
    count = 0
    average = 0
    while True:
        new_num = yield average
        if new_num is None:
            break
        count += 1
        total += new_num
        average = total/count
    return total, count, average


# 委托生成器
def proxy_gen():
    """
    如果把while True去掉的话 在调用方send(None)之后 子生成器结束后会return结果到委托生成器的total,count,average这三个变量上，
    但是在调用方send(None)后并没有接收到yield返回的值，它会触发StopIteration 。
    while True 使得委托生成器再进入到新的子生成器中，并在new_num = yield average处暂停返回给调用方。
    这样就不会引发异常，其实不用while True 在委托生成器的最后添加一个yield也是可以的
    """
    while True:
        yield from average_gen()


# 调用方
def main():
    calc_average = proxy_gen()
    next(calc_average)            # 预激下生成器
    print(calc_average.send(10))  # 打印：10.0
    print(calc_average.send(20))  # 打印：15.0
    print(calc_average.send(30))  # 打印：20.0


if __name__ == '__main__':
    # NOTE 首次 Next 在业务上可以理解为 预激, 首次调用 send(None) 也是预激，
    #      send 和 next 都会使生成器产出并让步
    #      预激：只执行 yield 右侧后就产出并让步，
    #      预激后：先执行 yield 左侧，直到碰到 yield，执行其右侧后产出并让步

    # g = yield0(15)
    # print(g)
    # print("--------------")
    # print(g.send(None))       # can't send non-None value to a just-started generator
    # print("--------------")
    # print(g.send(-1))
    # print("--------------")
    # print(next(g))
    # print("--------------")
    # print(g.send(-2))
    # print("--------------")
    # print(next(g))

    gF = yieldFrom(24)
    print("预激", next(gF))
    print("-1", gF.send(-1))
    print("-2", gF.send(-2))
    print("-3", gF.send(-3))
    print("-4", gF.send(-4))
    print("-5", gF.send(-5))

    """
    委托生成器的作用是：在调用方与子生成器之间建立一个双向通道。
    双向通道: 调用方可以通过send()直接发送消息给子生成器，而子生成器yield的值，也是直接返回给调用方。

    你可能会经常看到有些代码，还可以在yield from前面看到可以赋值。这是什么用法？
    你可能会以为，子生成器yield回来的值，被委托生成器给拦截了。你可以亲自写个demo运行试验一下，并不是你想的那样。
    因为我们之前说了，委托生成器，只起一个桥梁作用，它建立的是一个双向通道，它并没有权利也没有办法，对子生成器yield回来的内容做拦截。
    
    学到这里，我相信你肯定要问，既然委托生成器，起到的只是一个双向通道的作用，我还需要委托生成器做什么？
    我调用方直接调用子生成器不就好啦？   >>>>>>>>>>>>>>>> 因为它可以帮我们处理异常 <<<<<<<<<<<<<<<<<<<

    yield from帮我们做了很多的异常处理，而且全面，而这些如果我们要自己去实现的话，一个是编写代码难度增加，写出来的代码可读性极差，
    这些我们就不说了，最主要的是很可能有遗漏，只要哪个异常没考虑到，都有可能导致程序崩溃什么的
    """


def yieldFromHelp(EXPR):
    """
    一些说明
    _i：子生成器，同时也是一个迭代器
    _y：子生成器生产的值
    _r：yield from 表达式最终的值
    _s：调用方通过send()发送的值
    _e：异常对象
    """
    _i: Generator = iter(EXPR)
    try:
        _y = next(_i)
    except StopIteration as _e:
        _r = _e.value

    else:
        while 1:
            try:
                _s = yield _y
            except GeneratorExit as _e:
                try:
                    _m = _i.close
                except AttributeError:
                    pass
                else:
                    _m()
                raise _e
            except BaseException as _e:
                _x = sys.exc_info()
                try:
                    _m = _i.throw
                except AttributeError:
                    raise _e
                else:
                    try:
                        _y = _m(*_x)
                    except StopIteration as _e:
                        _r = _e.value
                        break
            else:
                try:
                    if _s is None:
                        _y = next(_i)
                    else:
                        _y = _i.send(_s)
                except StopIteration as _e:
                    _r = _e.value
                    break
    # RESULT = _r

    """
    1- 迭代器（即可指子生成器）产生的值直接返还给调用者
    2- 任何使用send()方法发给委派生产器（即外部生产器）的值被直接传递给迭代器。如果send值是None，则调用迭代器next()方法；
       如果不为None，则调用迭代器的send()方法。如果对迭代器的调用产生StopIteration异常，委派生产器恢复继续执行yield from后面的语句；
       若迭代器产生其他任何异常，则都传递给委派生产器。
    3- 子生成器可能只是一个迭代器，并不是一个作为协程的生成器，所以它不支持.throw()和.close()方法,即可能会产生AttributeError 异常。
    4- 除了GeneratorExit 异常外的其他抛给委派生产器的异常，将会被传递到迭代器的throw()方法。
       如果迭代器throw()调用产生了StopIteration异常，委派生产器恢复并继续执行，其他异常则传递给委派生产器。
    5- 如果GeneratorExit异常被抛给委派生产器，或者委派生产器的close()方法被调用，如果迭代器有close()的话也将被调用。
       如果close()调用产生异常，异常将传递给委派生产器。否则，委派生产器将抛出GeneratorExit 异常。
    6- 当迭代器结束并抛出异常时，yield from表达式的值是其StopIteration 异常中的第一个参数。
    7- 一个生成器中的return expr语句将会从生成器退出并抛出 StopIteration(expr)异常。
    """

    # TODO yield 或者 yield from 确实可以改变程序的执行流程，但是还是顺序执行的，这似乎与函数的直接调用没什么大的区别 ？？？
    # TODO 那实际的协程 是如何实现并行化的 ？？？？
