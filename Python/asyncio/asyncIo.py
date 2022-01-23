# -*- encoding: utf-8 -*-
import time
import asyncio
from asyncio.futures import Future
from collections.abc import Generator, Coroutine


"""
    yield 和 yield from 可以实现程序执行的流转，但是在我们的示例程序中，所有程序都是按照我们的调用顺序执行的，
    这似乎并没有体现并发的概念，对于 IO 密集型任务的调度又是如何实现并发的 ？
    IO 处理过程中，主程序继续执行，IO 处理过程中不会占用 CPU（线程） 资源吗 ？
"""


async def helloAsync(name):
    """
    将一个函数标记为协程
    """
    print('Hello,', name)


@asyncio.coroutine
def helloCoroutine():
    """
    只要在一个生成器函数头部用上 @asyncio.coroutine 装饰器就能将这个函数对象，【标记】为协程对象。
    注意这里是【标记】，划重点。实际上，它的本质还是一个生成器。标记后，它实际上已经可以当成协程使用。
    """
    # 异步调用asyncio.sleep(1):
    yield from asyncio.sleep(1)


def helloWorld():
    coroutine = helloAsync("World")               # 生成协程对象，并不会运行函数内的代码
    print(isinstance(coroutine, Coroutine))       # 检查是否是协程 Coroutine 类型 True
    # sys:1: RuntimeWarning: coroutine 'hello' was never awaited  vt. 等候，等待；期待

    coroutine = helloCoroutine()
    print(isinstance(coroutine, Generator), isinstance(coroutine, Coroutine))  # True, False
    # RuntimeWarning: Enable tracemalloc to get the object allocation traceback
# ----------------------------------------------------------------------------------------------------------------------


"""
event_loop 事件循环：
    程序开启一个无限的循环，程序员会把一些函数（协程）注册到事件循环上。当满足事件发生的时候，调用相应的协程函数。
coroutine 协程：
    协程对象，指一个使用async关键字定义的函数，它的调用不会立即执行函数，而是会返回一个协程对象。协程对象需要注册到事件循环，由事件循环调用。
future 对象：
    代表将来执行或没有执行的任务的结果。它和task上没有本质的区别
task 任务：
    一个协程对象就是一个原生可以挂起的函数，任务则是对协程进一步封装，其中包含任务的各种状态。
    Task 对象是 Future 的子类，它将 coroutine 和 Future 联系在一起，将 coroutine 封装成一个 Future 对象。
async/await 关键字：
    python3.5 用于定义协程的关键字，async定义一个协程，await用于挂起阻塞的异步调用接口。其作用在一定程度上类似于yield。

    1- 定义/创建协程对象
    2- 将协程转为task任务
    3- 定义事件循环对象容器
    4- 将task任务扔进事件循环对象中触发
"""


def asyncioTest():
    async def helloLL(name):
        print('Hello,', name)
    ll = helloLL("World")            # 定义协程对象
    loop = asyncio.get_event_loop()  # 定义事件循环对象容器
    # task = asyncio.ensure_future(ll)
    task = loop.create_task(ll)      # 将协程转为task任务
    loop.run_until_complete(task)    # 将task任务扔进事件循环对象中并触发

    """
        await用于挂起阻塞的异步调用接口。其作用在一定程度上类似于yield。
        注意这里是，一定程度上，意思是效果上一样（都能实现暂停的效果），但是功能上却不兼容。
        就是你不能在生成器中使用await，也不能在async 定义的协程中使用yield
        
        async 中不能使用yield， yield from 后面可接可迭代对象，也可接future对象或者协程对象；
        await 不能在未使用async修饰的函数中使用 ， await 后面必须要接 future对象或者协程对象；
        
        asyncio.sleep(n)，是asyncio自带的工具函数，可以模拟IO阻塞，并返回一个协程对象。
    """
    func = asyncio.sleep(2)
    print(isinstance(func, Future), isinstance(func, Coroutine), isinstance(task, Future))  # False  True

    # await 接协程
    async def hello1(name):
        await asyncio.sleep(2)
        print("Hello", name)

    # yield 接协程
    def hello2(name):  # async 修饰会报错
        yield from asyncio.sleep(2)
        print("Hello", name)

    # await 接future
    async def hello3(name):
        await asyncio.ensure_future(asyncio.sleep(2))
        print("Hello", name)

    # yield 接future
    def hello4(name):  # async 修饰会报错
        yield from asyncio.ensure_future(asyncio.sleep(2))
        print("Hello", name)

    hello1("hello1")  # NOTE 构建一个协程对象
    hello2("hello2")
    hello3("hello3")
    hello4("hello4")

    """
    异步IO的实现原理，就是在IO高的地方挂起，等IO结束后，再继续执行。
    在绝大部分时候，我们后续的代码的执行是需要依赖IO的返回值的，这就要用到回调了。
    回调的实现，有两种，一种是绝大部分程序员喜欢的，利用的同步编程实现的回调。
    这就要求我们要能够有办法取得协程的await的返回值。
    """
    async def _sleep(x):
        time.sleep(x)
        return '暂停了{}秒！'.format(x)
    coroutine = _sleep(2)
    loop = asyncio.get_event_loop()
    task = asyncio.ensure_future(coroutine)
    loop.run_until_complete(task)  # task.result() 可以取得返回结果
    print('返回结果：{}'.format(task.result()))

    """
    还有一种是通过asyncio自带的添加回调函数功能来实现。
    """
    async def _sleep(x):
        time.sleep(x)
        return '暂停了{}秒！'.format(x)

    def callback(future):
        print('这里是回调函数，获取返回结果是：', future.result())
    coroutine = _sleep(2)
    loop = asyncio.get_event_loop()
    task = asyncio.ensure_future(coroutine)
    task.add_done_callback(callback)  # 添加回调函数
    loop.run_until_complete(task)


if __name__ == '__main__':
    asyncioTest()
