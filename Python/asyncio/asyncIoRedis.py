# -*- encoding: utf-8 -*-

import time
from asyncio import AbstractEventLoop
from concurrent.futures import Future

import redis
import asyncio
from queue import Queue
from threading import Thread


# 将协程态添加到事件循环中
class EventLoopCoroutine(object):

    @staticmethod
    def start_loop(loop: AbstractEventLoop):
        # 一个在后台永远运行的事件循环
        # Equivalent to calling get_event_loop_policy().set_event_loop(loop).
        asyncio.set_event_loop(loop)
        loop.run_forever()
        print("after run_forever ...")

    @staticmethod
    def sync_do_sleep(x, que, mesg=""):
        print('More work {}'.format(x))
        time.sleep(x)     # 同步阻塞，loop 会等到该函数执行完毕才会去调用其它函数
        que.put(mesg)
        print('Finished more work {}'.format(x))

    @staticmethod
    async def async_do_sleep(x, que, mesg="ok"):
        # THINK await 是等待一个函数的执行直到返回一个 Future 对象, 包含（__await__ 的对象）
        # THINK async 标记一个函数为协程函数，在事件循环中执行该函数是要捕获阻塞式的IO 操作 ？？
        # THINK await 是主动等待 ？？？？？
        print('Waiting {}'.format(x))
        await asyncio.sleep(x)
        que.put(mesg)
        print('Done after {}s'.format(x))
        return f"-----> {x} --->"

    @staticmethod
    def get_redis():
        connection_pool = redis.ConnectionPool(host='127.0.0.1', db=0, password="holyzing")
        return redis.Redis(connection_pool=connection_pool)

    @staticmethod
    def consumer(rcon, que, event_loop):
        while True:
            task = rcon.rpop("queue")
            print(task)
            if not task:
                time.sleep(1)
                continue
            asyncio.run_coroutine_threadsafe(EventLoopCoroutine.async_do_sleep(int(task), que), event_loop)

    @staticmethod
    def coroutineThread(way="async"):
        """
        对于并发任务，通常是用生产者消费着模型，对队列的处理可以使用类似master-worker的方式，
        master主要为用户获取队列的msg，worker用户处理消息。
        为了简单起见，并且协程更适合单线程的方式，主线程用来监听队列，子线程用于处理队列。
        使用redis的队列。主线程中有一个是无限循环，用户消费队列。

        输入key=queue，value=5,3,1的消息

        loop_thread：单独的线程，运行着一个事件对象容器，用于实时接收新任务。
        consumer_thread：单独的线程，实时接收来自Redis的消息队列，并实时往事件对象容器中添加新任务。
        """
        queue = Queue()
        print(time.ctime())
        new_loop = asyncio.new_event_loop()

        # 定义一个新线程，并传入一个事件循环对象并启动，这样当前线程被阻塞
        # NOTE 可以设置为守护线程，在主线程结束的时候，子线程的执行也结束
        t = Thread(target=EventLoopCoroutine.start_loop, name='my', args=(new_loop,))
        t.setDaemon(True)
        t.start()

        if way != "redis":
            if way == "sync":      # 主线程中动态注册函数，不会阻塞主线程，新线程中按顺序执行注册的方法
                new_loop.call_soon_threadsafe(EventLoopCoroutine.sync_do_sleep, 6, queue, "第一个")
                new_loop.call_soon_threadsafe(EventLoopCoroutine.sync_do_sleep, 3, queue, "第二个")
                # THINK call_soon_threadsafe 不能加入协程对象 ???
            elif way == "async":   # 主线程中动态注册函数，不会阻塞主线程，新线程中并发执行注册的方法
                asyncio.run_coroutine_threadsafe(EventLoopCoroutine.async_do_sleep(
                    6, queue, "第一个"), new_loop).add_done_callback(lambda x: print(x.result()))
                fu: Future[str] = asyncio.run_coroutine_threadsafe(EventLoopCoroutine.async_do_sleep(
                    3, queue, "第二个"), new_loop)
                time.sleep(10)
                print("-------------------------------》 主线程睡眠结束---->")
                # THINK 添加回调函数可能会存在来不及调用的情况 ？？？？ 和协程的状态有关 ？？？
                fu.add_done_callback(lambda x: print(x.result()))
                # THINK (按添加顺序执行，遇到阻塞挂起执行下一个，执行完毕之后添加到后续执行容器，让另一个线程去继续循环执行)
                # THINK 也就是说一个事件循环中至少有两个线程在执行 ？？？？？？？？？？？？？？？？
            else:
                return "invalid async way"
        else:
            # 创建redis连接
            rcon = EventLoopCoroutine.get_redis()
            # 子线程：用于消费队列消息，并实时往事件对象容器中添加新任务
            consumer_thread = Thread(target=EventLoopCoroutine.consumer, args=(rcon, queue, new_loop))
            consumer_thread.setDaemon(True)
            consumer_thread.start()

        print("-------------- main thread --------------")
        while True:
            msg = queue.get()   # 默认是阻塞的
            print("{} 协程运行完..".format(msg))
            print(time.ctime())


if __name__ == '__main__':
    # EventLoopCoroutine.coroutineThread("sync")
    # EventLoopCoroutine.coroutineThread("async")
    EventLoopCoroutine.coroutineThread("redis")
    # ?????????????????????????????????????????????????????????????????????
    # THINK asyncio: 是如何检测到 IO 阻塞的，以及非阻塞的
    # ?????????????????????????????????????????????????????????????????????
