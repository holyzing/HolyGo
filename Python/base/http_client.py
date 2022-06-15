from flask import Flask

app = Flask(__name__)

@app.route("/")
def index():
    return "1"

if __name__ == '__main__':
    # tcpClient()
    app.run("192.168.1.2", 8080)

"""
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408019, 'ack_num': 0, 'data_offset': 8, 'flags': {'FIN': 0, 'SYN': 1, 'RST': 0, 'PSH': 0, 'ACK': 0, 'URG': 0}, 'win_size': 65535, 'tcp_checksum': 41099, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777453996, 'ack_num': 3899408020, 'data_offset': 8, 'flags': {'FIN': 0, 'SYN': 1, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 65535, 'tcp_checksum': 27356, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408020, 'ack_num': 1777453997, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 32218, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408020, 'ack_num': 1777453997, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 1, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 58980, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777453997, 'ack_num': 3899408498, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 31740, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777453997, 'ack_num': 3899408498, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 1, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 48158, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408498, 'ack_num': 1777454014, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 31723, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777454014, 'ack_num': 3899408498, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 1, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 48567, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408498, 'ack_num': 1777454150, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10232, 'tcp_checksum': 31588, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777454150, 'ack_num': 3899408498, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 1, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 19034, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408498, 'ack_num': 1777454151, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10232, 'tcp_checksum': 31587, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777454151, 'ack_num': 3899408498, 'data_offset': 5, 'flags': {'FIN': 1, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 31585, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408498, 'ack_num': 1777454152, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10232, 'tcp_checksum': 31586, 'urg_pointer': 0}
{'src_port': '192.168.1.2:2684', 'dst_port': '192.168.1.2:8080', 'seq_num': 3899408498, 'ack_num': 1777454152, 'data_offset': 5, 'flags': {'FIN': 1, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10232, 'tcp_checksum': 31585, 'urg_pointer': 0}
{'src_port': '192.168.1.2:8080', 'dst_port': '192.168.1.2:2684', 'seq_num': 1777454152, 'ack_num': 3899408499, 'data_offset': 5, 'flags': {'FIN': 0, 'SYN': 0, 'RST': 0, 'PSH': 0, 'ACK': 1, 'URG': 0}, 'win_size': 10233, 'tcp_checksum': 31584, 'urg_pointer': 0}

TODO Renova 中socket存在的问题，以及如何优化它， 网络编程的应用实战
"""

