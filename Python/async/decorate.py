# -*- encoding: utf-8 -*-


from re import A


def outer(*args, **kwargs):
    print(args, kwargs)

    def inner(func):
        res = func(*args, **kwargs)
        return res

    return inner



# -----------------------------------------------------
def outer(func):
    print(func)

    def inner(*args, **kwargs):
        print(args, kwargs)
        # 被装饰数作为内部函数,被装饰函数是没有参数的
        # res = func(*args, **kwargs)
        res = func()
        return res

    return inner

@outer
def ugly2():
    print("ll")


# 函数也是对象, 修改了 __name__, 因为返回的那个函数是 inner,只不过是用ugly2指向他而已,
# 为了不改变原函数的一些 "签名属性", 所以在 返回 inner之前可以把原始函数的一些 "签名属性"
# 复制到返回的函数. python内置了该功能 @functools.wraps(func)


print(ugly2.__name__)

ugly2(1, 2, a=2)

print("----------------------------------------------------")

def outerWithArgs(arg):
    print("first", arg)
    def decorate(func):
        print("second")
        def inner(*args, **kwargs):
            res = func(*args, **kwargs)
            return res
        print(inner.__name__)
        return inner
    print(decorate.__name__)
    return decorate


@outerWithArgs('ll')
def ugly():
    pass

print(ugly.__name__)


# 注意理解@ 他是一个特殊的操作符号 它的作用就是有 判断所接参数是否是一个有函数的参数(func),是否是一个类型(class)等

# first ll
# decorate
# second
# inner

# 从上述打印结果来看,发现outerWithArgs被执行了一次, decorate 被执行了一次 outerWithArgs(arg)(func)

# -----------------------------------------------------



def outer2():
    pass


@outer2
def ugly():
    pass


while True:
    pass
