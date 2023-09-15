import math
import time

# from joblib import Parallel, delayed


def my_fun_2p(i, j):
    """ We define a simple function with two parameters.
    """
    print(i, j, flush=True)
    return math.sqrt(i ** j)


print("准备开始执行并行任务 ...", flush=True)


def run():
    j_num = 6
    num = 10

    # start = time.time()
    # for i in range(num):
    #     for j in range(j_num):
    #         my_fun_2p(i, j)
    #
    # end = time.time()
    # print('{:.4f} s'.format(end - start))

    start = time.time()
    # n_jobs is the number of parallel jobs
    Parallel(n_jobs=7)(delayed(my_fun_2p)(i, j) for i in range(num) for j in range(j_num))

    end = time.time()
    print('{:.4f} s'.format(end - start))

    # UserWarning: A worker stopped while some jobs were given to the executor.
    # This can be caused by a too short worker timeout or by a memory leak.


def test_generator_expression():

    def show_list(*s):
        print(s)  # tuple
        print(type(s))

    def show_args(s):
        print(s)
        print(type(s))

    a = (str(i) for i in range(10))  # Generator
    show_list(a)
    show_args(str(i) for i in range(10))

def run2():
    import os
    print("Python Program")
    print(os.getenv("MYPATH"), flush=True)

if __name__ == '__main__':
    run2()


