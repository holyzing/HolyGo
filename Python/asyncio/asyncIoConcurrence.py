# -*- encoding: utf-8 -*-
import time
import random
import asyncio
import threading
from asyncio.futures import Future
from collections.abc import Coroutine  # Generator

"""
    将协程对象转为task任务对象
    定义一个事件循环对象容器用来存放task
    将task任务扔进事件循环对象中并触发

    协程中的并发
    协程中的嵌套
    协程中的状态
    gather与wait
"""


# 协程中的单任务
class SingleTaskTest(object):

    def __init__(self, name, count=0):
        self.name = name
        self.count = count

    async def coroutineFunc(self):
        while self.count < 20:
            await asyncio.sleep(2)    # 等待一个 Coroutine，并产生一个 Future
            self.count += 1
            print("------->>", "after last await", self.count)

    @asyncio.coroutine
    def generatorFunc(self):
        while self.count < 20:
            r = yield self.count
            self.count += 1
            print("------->>", "after last yield:", self.count, r)
        return self.count

    @asyncio.coroutine
    def generatorFunc2(self):
        while self.count < 20:
            # NOTE 去掉 @asyncio.coroutine 抛出异常
            # cannot 'yield from' a coroutine object in a non-coroutine generator
            r = yield from asyncio.sleep(2)
            self.count += 1
            print("------->>", "after last yield:", self.count, r)

    async def asyncIOFunc(self, sleepTime):
        print(self.name, "sleep begin")
        time.sleep(sleepTime)
        print(self.name, "sleep finished")
        return self.count

    @staticmethod
    def test():
        st = SingleTaskTest("ll", 10)
        c = st.coroutineFunc()        # 并不是直接调用一个函数, 而是产生一个协程对象
        print(c.send(None))
        # print(c.send(None))         # RuntimeError: await wasn't used with future
        print("---------------------------------------------------")
        c = st.generatorFunc()
        print(type(c))
        print(c.send(None))
        print(c.send(1))
        print("---------------------------------------------------")
        c = st.generatorFunc2()
        print(type(c))
        f: Future = c.send(None)
        f.set_result(2)
        print(f.result())
        print("---------------------------------------------------")
        c = st.asyncIOFunc(5)
        print(type(c))
        # c.send(None)
        print("---------------------------------------------------")
        loop = asyncio.get_event_loop()
        task = asyncio.ensure_future(c)
        loop.run_until_complete(task)


# 协程中的多任务
class MultiTaskTest(object):
    """
    asyncio实现并发，就需要多个协程来完成任务，每当有任务 "阻塞" 的时候就await，然后其他协程继续工作。
    """
    @staticmethod
    async def do_some_work(x: int):
        print('Waiting: ', x)
        await asyncio.sleep(x)
        print(f"IO 阻塞结束 {x} 继续执行")
        return 'Done after {}s'.format(x)

    @staticmethod
    def create_tasks():
        coroutine1 = MultiTaskTest.do_some_work(1)
        coroutine3 = MultiTaskTest.do_some_work(4)
        coroutine2 = MultiTaskTest.do_some_work(2)
        # 将协程转成task，并组成list
        task1 = asyncio.ensure_future(coroutine1)
        task3 = asyncio.ensure_future(coroutine3)
        task2 = asyncio.ensure_future(coroutine2)
        task3.add_done_callback(lambda x: print(f"回调函数3 处理 result:", x.result()))
        tasks = [task1, task2, task3]
        return tasks

    @staticmethod
    async def nestedCoroutines(method="wait"):
        # 使用async可以定义协程，协程用于耗时的io操作，我们也可以封装更多的io操作过程，
        # 这样就实现了嵌套的协程，即一个协程中await了另外一个协程，如此连接起来。
        # 【重点】：await 一个task列表（协程）
        # dones：表示已经完成的任务
        # pendings：表示悬而未决的任务

        # THINK 嵌套本身在业务逻辑上是没有什么实际意义的， 和 yield from 一样主要是为了作为一个桥梁，传递参数，处理其它异常而用 ？？？
        #       只是把创建协程对象，转换task任务，封装成在一个协程函数里而已。外部的协程，嵌套了一个内部的协程。
        #       而 “非嵌套”的单层线程运行，在 asyncio.wait() 源码中也是套了一层协程 _wait

        if method == "wait":
            dones, pendings = await asyncio.wait(MultiTaskTest.create_tasks())
            print("主协程 wait 中 所有子协程结束 ！")
            print(len(pendings))
            for task in dones:
                print('Task ret: ', task.result())
        elif method == "gather":
            results = await asyncio.gather(*MultiTaskTest.create_tasks())
            print("主协程 gather 中 所有子协程结束 ！")
            for result in results:
                print('Task ret: ', result)

    @staticmethod
    def run_tasks(tasks, method="wait"):
        loop = asyncio.get_event_loop()
        if isinstance(tasks, list):
            if method == "wait":
                loop.run_until_complete(asyncio.wait(tasks))
                print("主线程 wait 中所有协程结束 ！")
            elif method == "gather":
                loop.run_until_complete(asyncio.gather(*tasks))
                print("主线程 gather 中所有协程结束 ！")
            for task in tasks:
                print('Task ret: ', task.result())
        # NOTE 协程中嵌套一个多任务
        # elif callable(tasks):
        elif isinstance(tasks, Coroutine):
            loop.run_until_complete(tasks)
            print("主线程中所有协程结束 ！")

    @staticmethod
    def waitAndGatherDiff():
        loop = asyncio.get_event_loop()
        # group1 = asyncio.gather(*MultiTaskTest.create_tasks())
        # group2 = asyncio.gather(*MultiTaskTest.create_tasks())
        # group3 = asyncio.gather(*MultiTaskTest.create_tasks())
        # loop.run_until_complete(asyncio.gather(group1, group2, group3))

        print("主线程必须等待所有协程结束")

        # Wait 有控制功能
        async def coro(tag):
            t = random.uniform(0.5, 5)
            await asyncio.sleep(t)
            print(tag, "睡了", t)

        tasks = [coro(i) for i in range(1, 11)]
        # 最开始的肯定是第一个先执行，直到遇到 IO 则挂起执行第二个

        # 【控制运行任务数】：运行第一个任务就返回
        # FIRST_COMPLETED ：第一个任务完全返回
        # FIRST_EXCEPTION：产生第一个异常返回
        # ALL_COMPLETED：所有任务完成返回 （默认选项）
        dones, pendings = loop.run_until_complete(
            asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED))
        print("第一次完成的任务数:", len(dones))

        # 【控制时间】：运行一秒后，就返回
        dones2, pendings2 = loop.run_until_complete(
            asyncio.wait(pendings, timeout=3))
        print("第二次完成的任务数:", len(dones2))

        # 【默认】：所有任务完成后返回
        dones3, pendings3 = loop.run_until_complete(asyncio.wait(pendings2))

        print("第三次完成的任务数:", len(dones3))

        loop.close()


def coroutineStateTest():
    """
    Future对象，或者Task任务 的状态
    Pending：创建future，还未执行
    Running：事件循环正在调用执行任务
    Done：任务执行完毕
    Cancelled：Task被取消后的状态
    """
    async def hello():
        print("Running in the loop...")
        flag = 0
        while flag < 1000000:
            with open("/home/holyzing/xtalpi/My/_03_Scala/Scala/xtalpi/tmp/test.txt", "a") as f:
                f.write("\n------")
            flag += 1
        print("Stop the loop")

    coroutine = hello()
    loop = asyncio.get_event_loop()
    # task = asyncio.ensure_future(coroutine)
    task = loop.create_task(coroutine)

    print(task)         # Pending：未执行状态
    try:
        t1 = threading.Thread(target=loop.run_until_complete, args=(task,))
        # t1.daemon = True      # TODO 多线程需要好好理解一下
        t1.start()
        time.sleep(4)
        print("----------------------------------------------》")
        print(task)     # Pending：运行中状态
        print("加入到当前主线程 ...")
        t1.join()
        print("等待 t1 执行完成之后，才会notify 当前线程 ！")
    except KeyboardInterrupt as e:
        print(e)
        task.cancel()   # Cacelled：取消任务
        print(task)
    finally:
        print(task)

    """
    IMPORTANT NOTE :
    你对Pending这个词的意思有误解，pending在这里的意思是悬而未决的意思。
    所以不存在你说的有pending和pending running这两个状态之分，在task没结束之前，
    他的状态都是pending，我发现中国人学编程最大的困难是英语差

    关于asyncio.wait接收参数的方式一节中，对coro不转为task对象传入的方式其实在文档中已经有说明“直接向 wait() 传入协程对象的方式已弃用”
    https://docs.python.org/zh-cn/3.7/library/asyncio-task.html?highlight=wait#asyncio.wait
    """


if __name__ == '__main__':
    """
     NOTE asyncio 是如何 捕获到 IO 阻塞开始，以及阻塞结束的 ？
    """
    # SingleTaskTest.test()
    # MultiTaskTest.run_tasks(MultiTaskTest.create_tasks(), "wait")
    # MultiTaskTest.run_tasks(MultiTaskTest.create_tasks(), "gather")
    # MultiTaskTest.run_tasks(MultiTaskTest.nestedCoroutines("wait"))
    # MultiTaskTest.run_tasks(MultiTaskTest.nestedCoroutines("gather"))
    # coroutineStateTest()
    MultiTaskTest.waitAndGatherDiff()



