# -*- encoding: utf-8 -*-
import time
import asyncio
import requests
"""
 IO 多路复用：

    |应用层 （socket Client）
    |   socket 是介于 应用层与运输层之间的抽象接口，它抽象了应用和TCP协议镞之间的交流，对外简化了它们之间通信的步骤
    |运输层（TCP/UDP 协议层）
    |网络层（ICMP IP ICMP）
    |链路层（ARP 硬件接口 RARP）

1 - 使用 cgroups 和 namespace 实现进程间隔离
2 - 进程间通信

应用层
表示层
会话层
传输层
网络层
链路层
物理层

"""

"""
协程的原理：
1- python 在语法层面 限制 yield 后必须接着生成器（非协程可迭代的对象），await 后必须接着协程（Future和Task）。
2- yield 和 yield from 都是 由用户实现多个生成器上下文的切换，而 await则是由解释器在一定环境下（eventloop）实现多个协程间上下文的切换。

3- 解释器会检测到 IO操作，并挂起当前执行，去执行其它协程 ？？？？？
4- 解释器会检测到 IO执行完成，并继续执行上下文
5- asyncio 并不会自动捕获 httpio 然后挂起去执行其它函数
6- 需要使用 with 并借助微线程实现异步并发执行

"""


async def asyncReq(name):
    r = requests.get(f"http://10.0.7.127:6060/admin/temp/test/httpio/{name}")
    return r


async def httpIO(name):
    # print("sleep1 Start  ----------------------->", name)
    # time.sleep(5)
    # print("sleep1 end    ----------------------->", name)
    print("httpIO Start ----------------------->", name)
    # r = requests.get(f"http://10.0.7.127:6060/admin/temp/test/httpio/{name}")
    r = asyncio.ensure_future(asyncReq(name))
    r.add_done_callback(callback)
    await r
    print("httpIO end ----------------------->", name)
    # print("sleep2 Start  ----------------------->", name)
    # time.sleep(5)
    # print("sleep2 end    ----------------------->", name)


def callback(r):
    print(r.result().json())


async def testDetectIO(name):
    print("outer Start ----------------------->", name)
    hi1 = httpIO("A")
    # hi1 = asyncio.ensure_future(hi1)
    # hi1.add_done_callback(callback)
    hi2 = httpIO("B")
    # hi2 = asyncio.ensure_future(hi2)
    # hi2.add_done_callback(callback)
    print("await Start ----------------------->", name)
    dones, pendings = await asyncio.wait([hi1, hi2])
    print("await end   ----------------------->", name)

    # print(len(pendings))
    for task in dones:
        print('Task ret: ', task.result())
    print("outer end   ----------------------->", name)


if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(testDetectIO("main"))
