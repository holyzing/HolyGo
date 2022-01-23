# -*- encoding: utf-8 -*-

"""
    需要借助collections.abc这个模块（Python2没有），使用isinstance()来
    判别一个对象是否是可迭代的（Iterable），是否是迭代器（Iterator），是否是生成器（Generator）

    可迭代对象，是其内部实现了，__iter__ 这个魔术方法。
    可以通过，dir()方法来查看是否有__iter__来判断一个变量是否是可迭代的。

    迭代器, 对比可迭代对象，迭代器其实就只是多了一个函数而已。就是 __next__()，我们可以不再使用for循环来间断获取元素值。
    而可以直接使用next()方法来实现。迭代器，是在可迭代的基础上实现的。要创建一个迭代器，我们首先，得有一个可迭代对象。

    引入生成器，是为了实现一个在计算下一个值时不需要浪费空间的结构，生成器，
    则是在迭代器的基础上（可以用for循环，可以使用next()），再实现了yield

    yield 是什么东西呢，它相当于我们函数里的return。在每次next()，或者for遍历的时候，
    都会yield这里将新的值返回回去，并在这里阻塞，等待下一次的调用。
    正是由于这个机制，才使用生成器在Python编程中大放异彩。实现节省内存，实现异步编程。

    可迭代对象和迭代器，是将所有的值都生成存放在内存中，而生成器则是需要元素才临时生成，节省时间，节省空间
    由于生成器并不是一次生成所有元素，而是一次一次的执行返回，那么如何刺激生成器执行(或者说激活)呢？
"""

import collections
from collections.abc import Iterable, Iterator, Generator


def iterableTest():
    # 字符串
    astr = "XiaoMing"
    print("字符串：{}".format(astr), isinstance(astr, Iterable),
          isinstance(astr, Iterator), isinstance(astr, Generator))

    # 列表
    alist = [21, 23, 32, 19]
    print("列表：{}".format(alist), isinstance(alist, Iterable),
          isinstance(alist, Iterator), isinstance(alist, Generator))

    # 字典
    adict = {"name": "小明", "gender": "男", "age": 18}
    print("字典：{}".format(adict), isinstance(adict, Iterable),
          isinstance(adict, Iterator), isinstance(adict, Generator))

    # deque
    deque = collections.deque('abcdefg')
    print("deque：{}".format(deque), isinstance(deque, Iterable),
          isinstance(deque, Iterator), isinstance(deque, Generator))

    # generator
    gen = (i for i in range(3))
    print("gen：{}".format(gen), isinstance(gen, Iterable),
          isinstance(gen, Iterator), isinstance(gen, Generator))


def iteratorTest():

    class MyList(object):           # 定义可迭代对象类
        def __init__(self, num):
            self.end = num          # 上边界

        def __iter__(self):         # 返回一个实现了__iter__和__next__的迭代器类的实例
            return MyListIterator(self.end)

    class MyListIterator(object):   # 定义迭代器类
        def __init__(self, end):
            self.data = end         # 上边界
            self.start = 0

        def __iter__(self):         # 返回该对象的迭代器类的实例；因为自己就是迭代器，所以返回self
            return self

        def __next__(self):         # 迭代器类必须实现的方法，若是Python2则是next()函数
            while self.start < self.data:
                self.start += 1
                return self.start - 1
            raise StopIteration

    my_list = MyList(5)                       # 得到一个可迭代对象
    print(isinstance(my_list, Iterable))      # True
    print(isinstance(my_list, Iterator))      # False
    for i in my_list:                         # 迭代
        print(i)

    my_iterator = iter(my_list)               # 得到一个迭代器
    print(isinstance(my_iterator, Iterable))  # True
    print(isinstance(my_iterator, Iterator))  # True

    print(next(my_iterator))                  # 迭代
    print(next(my_iterator))
    print(next(my_iterator))
    print(next(my_iterator))
    print(next(my_iterator))

    aStr = 'abcd'           # 创建字符串，它是可迭代对象
    aIterator = iter(aStr)  # 通过iter()，将可迭代对象转换为一个迭代器
    print(isinstance(aIterator, Iterator))  # True
    next(aIterator)  # a


def generatorTest():
    # 使用列表生成式，注意不是[]，而是()
    L = (x * x for x in range(10))
    print(isinstance(L, Generator))  # True

    # 实现了yield的函数
    def mygen(n):
        now = 0
        while now < n:
            yield now
            now += 1

    mygen(1)


def prerequisiteTest():
    # 使用next()
    # 使用generator.send(None)
    pass


def generatorStateTest():
    """
    GEN_CREATED # 等待开始执行
    GEN_RUNNING # 解释器正在执行（只有在多线程应用中才能看到这个状态）
    GEN_SUSPENDED # 在yield表达式处暂停
    GEN_CLOSED # 执行结束
    """
    from inspect import getgeneratorstate

    def mygen(n):
        now = 0
        while now < n:
            yield now
            now += 1
    gen = mygen(2)
    print(getgeneratorstate(gen))

    # TODO GEN_RUNNING这个状态，使用多线程获取其状态即可

    print(next(gen))
    print(getgeneratorstate(gen))

    print(next(gen))
    gen.close()  # 手动关闭/结束生成器
    print(getgeneratorstate(gen))


def stopIterationTest():
    """
    在生成器工作过程中，若生成器不满足生成元素的条件，就会|应该 抛出异常（StopIteration）。
    通过列表生成式构建的生成器，其内部已经自动帮我们实现了抛出异常这一步
    所以我们在自己定义一个生成器的时候，我们也应该在不满足生成元素条件的时候，抛出异常。
    """
    def mygen(n):
        now = 0
        while now < n:
            yield now
            now += 1
        raise StopIteration

    gen = mygen(2)
    next(gen)
    next(gen)
    next(gen)


"""
    生成器为我们引入了暂停函数执行（yield）的功能。当有了暂停的功能之后，人们就想能不能在生成器暂停的时候向其发送一点东西
    （其实上面也有提及：send(None)）。这种向暂停的生成器发送信息的功能通过 PEP 342 进入 Python 2.5 中，
    并催生了 Python 中协程的诞生。根据 wikipedia 中的定义
    
        协程是为非抢占式多任务产生子程序的计算机程序组件，协程允许不同入口点在不同位置暂停或开始执行程序
    
    注意从本质上而言，协程并不属于语言中的概念，而是编程模型上的概念。
    协程和线程，有相似点，多个协程之间和线程一样，只会交叉串行执行；也有不同点，线程之间要频繁进行切换，加锁，解锁，从复杂度和效率来看，
    和协程相比，这确是一个痛点。协程通过使用 yield 暂停生成器，可以将程序的执行流程交给其他的子程序，从而实现不同子程序的之间的交替执行。
"""


if __name__ == '__main__':
    generatorStateTest()
