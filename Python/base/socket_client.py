# -*- encoding: utf-8 -*-


import socket
def tcpClient():
    host = "192.168.1.2"
    port = 8080

    s = socket.socket()
    s.connect((host, port))
    message = input("->")
    while True:
        data = s.recv(1024)
        print(data)
        print(data.decode())
        s.send(message)
        message = input("->")
    s.close()



