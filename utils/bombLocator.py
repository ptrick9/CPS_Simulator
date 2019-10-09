import re
import random


f = open('/home/simulator/git-simulator/movement/Firefly_1000_3.scb', 'r')

coord = re.compile(r'(-?\d+,-?\d+)')

choices = set()


for i, line in enumerate(f):
    if i > 500:
        x = re.findall(coord, line)
        for val in x:
            nums = val.split(',')
            nums = [int(y) for y in nums]
            if nums[0] > 0 and nums[1] > 0:
                choices.add(val)
print(len(choices))

choices = list(choices)
random.shuffle(choices)

chosen = choices[:40]

print(chosen)

chosen = [[int(y) for y in x.split(',')] for x in chosen]

print(chosen)

for i, val in enumerate(chosen):
    for j,v in enumerate(val):
        chosen[i][j] = v + random.choice([-4, -3, -2, -1, 0, 1, 2, 3, 4])

print(chosen)

s = 'commBomb = ["-commandBomb=%s -bombX=%d -bombY=%d" % (s[0], s[1], s[2]) for s in ['

for val in chosen:
    ss = '("true", %d, %d),' % (val[0], val[1])
    s += ss

s = s[:-1]
s += ']]'

print(s)