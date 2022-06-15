# -*- encoding: utf-8 -*-

import uuid
import json
import psutil
import socket
import struct
import threading

from requests import request


class ThreadSafeQueue(object):

    def __init__(self, max_size=0):
        self.queue = []
        self.max_size = max_size  # max_size为0表示无限大
        self.lock = threading.Lock()  # 互斥量
        self.condition = threading.Condition()  # 条件变量

    def size(self):
        """
        获取当前队列的大小
        :return: 队列长度
        """
        # 加锁
        self.lock.acquire()
        size = len(self.queue)
        self.lock.release()
        return size

    def put(self, item):
        """
        将单个元素放入队列
        :param item:
        :return:
        """
        # 队列已满 max_size为0表示无限大
        if self.max_size != 0 and self.size() >= self.max_size:
            return ThreadSafeException()

        # 加锁
        self.lock.acquire()
        self.queue.append(item)
        self.lock.release()
        self.condition.acquire()
        # 通知等待读取的线程
        self.condition.notify()
        self.condition.release()

        return item

    def batch_put(self, item_list):
        """
        批量添加元素
        :param item_list:
        :return:
        """
        if not isinstance(item_list, list):
            item_list = list(item_list)

        res = [self.put(item) for item in item_list]

        return res

    def pop(self, block=False, timeout=0):
        """
        从队列头部取出元素
        :param block: 是否阻塞线程
        :param timeout: 等待时间
        :return:
        """
        if self.size() == 0:
            if block:
                self.condition.acquire()
                self.condition.wait(timeout)
                self.condition.release()
            else:
                return None

        # 加锁
        self.lock.acquire()
        item = None
        if len(self.queue):
            item = self.queue.pop()
        self.lock.release()

        return item

    def get(self, index):
        """
        获取指定位置的元素
        :param index:
        :return:
        """
        if self.size() == 0 or index >= self.size():
            return None

        # 加锁
        self.lock.acquire()
        item = self.queue[index]
        self.lock.release()

        return item


class ThreadSafeException(Exception):
    pass


class ThreadProcess(threading.Thread):

    def __init__(self, task_queue, *args, **kwargs):
        """
        线程处理方法初始化
        :param task_queue:
        :param args:
        :param kwargs:
        """
        super(ThreadProcess, self).__init__(*args, **kwargs)
        self.dismiss_flag = threading.Event()  # 任务停止的标记
        self.task_queue = task_queue
        self.args = args
        self.kwargs = kwargs

    def run(self):
        """
        线程运行方法
        :return:
        """
        while True:
            # 线程停止标志设定则停止执行
            if self.dismiss_flag.is_set():
                break
            # task对象是否是Task的实例
            task = self.task_queue.pop()
            if not isinstance(task, Task):
                continue

            # print('task id:%d' % task.id)
            # print(type(task))
            result = task.callable(*task.args, **task.kwargs)
            # 如果是异步任务 设置返回结果
            if isinstance(task, AsyncTask):
                # print('set result:%d' % task.id)
                task.set_result(result)

    def __dismiss(self):
        self.dismiss_flag.set()

    def stop(self):
        """
        线程停止方法
        :return:
        """
        self.__dismiss()


class ThreadPool(object):

    def __init__(self, size=0):
        if not size:
            size = psutil.cpu_count() * 2
        print('size is %d' % size)
        self.size = size
        # 任务队列 存放待处理的任务
        self.task_queue = ThreadSafeQueue()
        # 线程池
        self.pool = ThreadSafeQueue(size)

        for i in range(self.size):
            print('put thread:%d' % i)
            self.pool.put(ThreadProcess(self.task_queue))

    def start(self):
        """
        开启线程池
        :return:
        """
        for i in range(self.pool.size()):
            print('start thread:%d' % i)
            thread = self.pool.get(i)
            thread.start()

    def join(self):
        """
        停止线程池
        :return:
        """
        for i in range(self.pool.size()):
            print('join thread:%d' % i)
            thread = self.pool.get(i)
            thread.stop()
        # 线程池未完全停止
        if not self.pool:
            thread = self.pool.pop()
            thread.join()

    def put(self, task):
        """
        将任务放入线程池
        :param task:
        :return:
        """
        if not isinstance(task, Task):
            raise TaskTypeException()
        res = self.task_queue.put(task)
        # print('put task %d' % task.id)

        return res

    def batch_put(self, task_list):
        """
        向线程池中批量添加任务
        :param task_list:
        :return:
        """
        if not isinstance(task_list, list):
            task_list = list(task_list)

        return [self.put(task) for task in task_list]

    def size(self):
        """
        获取线程池的大小
        :return:
        """

        return self.pool.size()


class TaskTypeException(Exception):
    pass


def TestStructUnpack():
    bin_str = b'ABCDEFGH'
    print(bin_str)
    print(bin_str.decode())
    res = struct.unpack('>8B', bin_str)
    print(res)
    res2 = struct.unpack('>4H', bin_str)
    print(res2)
    res3 = struct.unpack('>2L', bin_str)
    print(res3)
    res4 = struct.unpack('>8s', bin_str)
    print(res4)


class Task(object):

    def __init__(self, func, *args, **kwargs):
        """
        基本任务对象初始化
        :param func:
        :param args:
        :param kwargs:
        """
        self.id = uuid.uuid4()
        self.callable = func
        self.args = args
        self.kwargs = kwargs

    def __str__(self):
        return 'Task id: %s' % str(self.id)


class AsyncTask(Task):

    def __init__(self, func, *args, **kwargs):
        """
        异步任务对象初始化
        :param callable:
        """
        super(AsyncTask, self).__init__(func, *args, **kwargs)
        self.result = None
        self.condition = threading.Condition()

    def set_result(self, result):
        """
        设置返回结果
        :param result:
        :return:
        """
        # 加锁 通知被阻塞线程
        self.condition.acquire()
        self.result = result
        self.condition.notify()
        self.condition.release()

    def get_result(self):
        """
        获取返回结果 结果不存在时会阻塞当前线程
        :return:
        """
        # 加锁 等待任务处理结果写入
        self.condition.acquire()
        if not self.result:
            self.condition.wait()
        result = self.result
        self.condition.release()

        return result


class TransParser(object):

    IP_HEADER_LENGTH = 20  # IP报文头部的长度
    UDP_HEADER_LENGTH = 8  # UDP头部的长度
    TCP_HEADER_LENGTH = 20  # TCP头部的长度


class TCPParser(TransParser):

    @classmethod
    def parse_tcp_header(cls, tcp_header):
        """
        TCP报文格式
        1. 16位源端口号 16位目的端口号
        2. 32位***
        3. 32位确认号
        4. 4位数据偏移 6位保留字段 6位TCP标记 16位窗口
        5. 16位校验和 16位紧急指针
        :param tcp_header:
        :return:
        """
        line1 = struct.unpack('>HH', tcp_header[:4])
        src_port = line1[0]
        dst_port = line1[1]
        line2 = struct.unpack('>L', tcp_header[4:8])
        seq_num = line2[0]
        line3 = struct.unpack('>L', tcp_header[8:12])
        ack_num = line3[0]
        line4 = struct.unpack('>BBH', tcp_header[12:16])  # 先按照8位、8位、16位解析
        data_offset = line4[0] >> 4  # 第一个8位右移四位获取高四位
        flags = line4[1] & int(b'00111111', 2)  # 第二个八位与00111111进行与运算获取低六位
        FIN = flags & 1
        SYN = (flags >> 1) & 1
        RST = (flags >> 2) & 1
        PSH = (flags >> 3) & 1
        ACK = (flags >> 4) & 1
        URG = (flags >> 5) & 1
        win_size = line4[2]
        line5 = struct.unpack('>HH', tcp_header[16:20])
        tcp_checksum = line5[0]
        urg_pointer = line5[1]

        # 返回结果
        # src_port 源端口
        # dst_port 目的端口
        # seq_num ***
        # ack_num 确认号
        # data_offset 数据偏移量
        # flags 标志位
        #     FIN 结束位
        #     SYN 同步位
        #     RST 重启位
        #     PSH 推送位
        #     ACK 确认位
        #     URG 紧急位
        # win_size 窗口大小
        # tcp_checksum TCP校验和
        # urg_pointer 紧急指针
        return {
            'src_port': src_port,
            'dst_port': dst_port,
            'seq_num': seq_num,
            'ack_num': ack_num,
            'data_offset': data_offset,
            'flags': {
                'FIN': FIN,
                'SYN': SYN,
                'RST': RST,
                'PSH': PSH,
                'ACK': ACK,
                'URG': URG
            },
            'win_size': win_size,
            'tcp_checksum': tcp_checksum,
            'urg_pointer': urg_pointer
        }

    @classmethod
    def parser(cls, packet):

        return cls.parse_tcp_header(packet[cls.IP_HEADER_LENGTH:cls.IP_HEADER_LENGTH + cls.TCP_HEADER_LENGTH])


class IPParser(object):

    IP_HEADER_LENGTH = 20  # 报文前二十字节为ip头部

    @classmethod
    def parse_ip_header(cls, ip_header):
        """
        IP报文格式
        1. 4位IP-version 4位IP头长度 8位服务类型 16位报文总长度
        2. 16位标识符 3位标记位 13位片偏移 暂时不关注此行
        3. 8位TTL 8位协议 16位头部校验和
        4. 32位源IP地址
        5. 32位目的IP地址
        :param ip_header:
        :return:
        """
        line1 = struct.unpack('>BBH', ip_header[:4])  # 先按照8位、8位、16位解析
        ip_version = line1[0] >> 4  # 通过右移4位获取高四位
        # 报文头部长度的单位是32位 即四个字节
        iph_length = (line1[0] & 15) * 4  # 与1111与运算获取低四位
        packet_length = line1[2]
        line3 = struct.unpack('>BBH', ip_header[8: 12])
        TTL = line3[0]
        protocol = line3[1]
        iph_checksum = line3[2]
        line4 = struct.unpack('>4s', ip_header[12: 16])
        src_ip = socket.inet_ntoa(line4[0])
        line5 = struct.unpack('>4s', ip_header[16: 20])
        dst_ip = socket.inet_ntoa(line5[0])

        # 返回结果
        # ip_version ip版本
        # iph_length ip头部长度
        # packet_length 报文长度
        # TTL 报文寿命
        # protocol 协议号 1 ICMP协议 6 TCP协议 17 UDP协议
        # iph_checksum ip头部的校验和
        # src_ip 源ip
        # dst_ip 目的ip
        return {
            'ip_version': ip_version,
            'iph_length': iph_length,
            'packet_length': packet_length,
            'TTL': TTL,
            'protocol': protocol,
            'iph_checksum': iph_checksum,
            'src_ip': src_ip,
            'dst_ip': dst_ip
        }

    @classmethod
    def parse(cls, packet):
        ip_header = packet[:cls.IP_HEADER_LENGTH]

        return cls.parse_ip_header(ip_header)


class ServerProcessTask(AsyncTask):

    def __init__(self, packet, *args, **kwargs):
        """
        定义异步处理任务
        :param packet:
        :param args:
        :param kwargs:
        """
        super(ServerProcessTask, self).__init__(func=self.process, *args, **kwargs)
        self.packet = packet

    def process(self):
        """
        异步处理方法
        :return:
        """
        headers = {
            'network_header': None,
            'transport_header': None
        }

        ip_header = IPParser.parse(self.packet)
        headers['network_header'] = ip_header
        if ip_header['protocol'] == 6:
            headers['transport_header'] = TCPParser.parser(self.packet)

        return headers


class Server(object):

    # 在win下才支持混杂模式

    def __init__(self):
        # 创建socket 指明工作协议类型(IPv4) 套接字类型 工作具体的协议(IP协议)
        # win
        self.sock:socket.socket = socket.socket(socket.AF_INET, socket.SOCK_RAW, socket.IPPROTO_IP)

        # 设置自己的主机ip和端口
        self.ip = '192.168.1.2'
        self.port = 8080
        self.sock.bind((self.ip, self.port))

        # self.sock.listen(1)

        # 设置混杂模式 接受所有经过网卡设备的数据
        self.sock.ioctl(socket.SIO_RCVALL, socket.RCVALL_ON)


        # 初始化线程池
        self.pool = ThreadPool(10)
        self.pool.start()

    def loop_server(self):
        """
        循环读取网络数据
        :return:
        """
        while True:
            packet, addr = self.sock.recvfrom(65535)
            task = ServerProcessTask(packet)
            self.pool.put(task)
            result = task.get_result()
            if result["transport_header"] != None and (result["transport_header"]["src_port"] == 8080 or result["transport_header"]["dst_port"] == 8080):
                # result = json.dumps(result, indent=4)
                result["transport_header"]["src_port"] = result["network_header"]["src_ip"] + ":" + str(result["transport_header"]["src_port"])
                result["transport_header"]["dst_port"] = result["network_header"]["dst_ip"] + ":" + str(result["transport_header"]["dst_port"])

                print(result["transport_header"])

Server().loop_server()
