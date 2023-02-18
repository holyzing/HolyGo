'''
兜兜算式式
'''
# -*- encoding: utf-8 -*-
import random

s = set()
while True:
    a = random.randint(0, 20)
    b = random.randint(0, 20)
    if a - b < 0:
        continue
    s.add(f"{a}-{b}=")
    if a + b > 20:
        continue
    s.add(f"{a}+{b}=")
    if len(s) > 200:
        break

for i in s:
    print(i)
