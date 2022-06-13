# -*- encoding: utf-8 -*-

class Area(object):
    # 对未进行初始化的class 成员变量,进行访问会报错
    MyName: str  # 类型会影响解释的行为
    # MyName: str = "Area"


class State(Area):
    Age: str

def show(area: Area):
    print(area.MyName)

# Area.MyName = "Area"
show(Area())
show(State())
