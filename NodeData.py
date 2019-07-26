import re
import sys
from PIL import Image, ImageDraw, ImageColor
import random

class NodeData:
    recalTimes = []
    selfRecalTimes = []

    def __init__(self, id, e0, e1, e2, et1, et2, s0, s1, s2):
        self.id = id
        self.s0 = s0
        self.s1 = s1
        self.s2 = s2
        self.e0 = e0
        self.e1 = e1
        self.e2 = e2
        self.et1 = et1
        self.et2 = et2

def getNodeData():
    log = open("Log-nodeData.txt", 'r')
    if log.mode == 'r':
        lines = log.readlines()
    log.close()

    nodes = []

    for i, line in enumerate(lines):
        data = re.split(",", line)
        recalTime = []
        selfRecal = []
        for j, val in enumerate(data):
            if val[:2] == "ID":
                id = val[2:]
            elif val[:2] == "E0":
                e0 = val[2:]
            elif val[:2] == "E1":
                e1 = val[2:]
            elif val[:2] == "E2":
                e2 = val[2:]
            elif val[:3] == "ET1":
                et1 = val[3:]
            elif val[:3] == "ET2":
                et2 = val[3:]
            elif val[:2] == "S0":
                s0 = val[2:]
            elif val[:2] == "S1":
                s1 = val[2:]
            elif val[:2] == "S2":
                s2 = val[2:]
            elif val[:2] == "RT":
                recalStr = re.sub("[\[\]]", "", val[2:])
                recalStr = re.split(" ", recalStr)
                for num in recalStr:
                    recalTime.append(int(num))
            elif val[:3] == "SRT":
                selfRecalStr = re.sub("[\[\]]", "", val[3:])
                selfRecalStr = re.split(" ", selfRecalStr)
                for num in selfRecalStr:
                    selfRecal.append(int(num))
        n = NodeData(id, e0, e1, e2, et1, et2, s0, s1, s2)
        n.recalTimes = recalTime
        n.selfRecalTimes = selfRecal
        nodes.append(n)
    return nodes

#for node in nodes:
#    print(node.id, node.recalTimes, node.selfRecalTimes)
