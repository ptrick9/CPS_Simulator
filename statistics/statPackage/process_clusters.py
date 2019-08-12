import os
import pandas as pd
import csv
import numpy as np
from zipfile import *
from io import StringIO

def isFloat(s):
    try:
        float(s)
        return True
    except ValueError:
        #print("Value error: ",s)
        return False

class Cluster:
    def __init__(self, line):
        self.line = line
        self.ch = ""
        self.cm = ""

    def process(self):
        [self.ch,self.cm] = self.line.split('-')
        self.ch = (self.ch).split('H')[1]
        self.cm = (self.cm).split(' ')
        new_set = set()
        for i in range(0,len(self.cm)):
            self.cm[i]
            self.cm[i] = self.cm[i].replace("M", "")
            self.cm[i] = self.cm[i].replace("[", "")
            self.cm[i] = self.cm[i].replace(",", "")
            self.cm[i] = self.cm[i].replace("]", "")
            new_set.add(self.cm[i])
        self.cm = new_set

def getBatteryStats(basename):

    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-batteryusage.csv")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    with filename as csv_file:
        data = list(csv.reader(csv_file, delimiter=','))
        iters = len(data)
        num_iters = len(data)-1
        num_nodes = len(data[0])-1
        batteryLevels = [x[:] for x in [[0.0] * num_nodes] * num_iters]
        #print(num_nodes)

        for i in range(0,num_iters-1):
            for j in range(1,num_nodes):

                if(isFloat(data[i][j])):
                    try:
                        batteryLevels[i][j-1] = float(data[i][j])
                    except IndexError:
                        #print(i,j)
                        pass

        b = [x[:] for x in [[0.0] * (num_nodes-2)] * (num_iters-1)]
        for i in range(1,num_iters-1):
            b[i-1] = sorted(batteryLevels[i][:num_iters])

        minimum = [0 for i in range(0,iters)]
        maximum = [0 for i in range(0,iters)]
        lowerq = [0 for i in range(0,iters)]
        upperq = [0 for i in range(0,iters)]
        median = [0 for i in range(0,iters)]
        end = len(b)

        for i in range(0,num_iters-1):
            minimum[i] = b[i][0]
            maximum[i] = b[i][end]
            lowerq[i] = b[i][int((end+2)/4)]
            upperq[i] = b[i][int(3*(end+2)/4)]
            median[i] = b[i][int((end+2)/2)]
        
    return (minimum, lowerq, median, upperq, maximum)

def clusterHeadsPerIter(basename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-clusters.txt")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    #filename = basename + "-clusters.txt"
    with filename as csv_file:
        data = list(csv.reader(csv_file, delimiter=','))
        iters = len(data)
        numCH = [None]*iters
        count = 0
        for row in data:
            numCH[count] = len(row)
            count = count+1
            for i in range(0,len(row)):
                row[i] = int(row[i])
    t = range(iters)
#    fig = figure(num=None, figsize=(8, 6), dpi=90, facecolor='w', edgecolor='k')
#    plt.plot(t, numCH, 'ro')
#    plt.title("Total Cluster Heads at each Iteration")
#    plt.xlabel("Iteration")
#    plt.ylabel("# of Cluster Heads")
#    plt.show()
    
    return numCH

def getIters(basename, filename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, filename)
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    with filename as csv_file:
        data = list(csv.reader(csv_file, delimiter=','))
    return len(data)

def clusterMessagesPerIter(basename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-cluster_messages.txt")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    #filename = basename + "-cluster_messages.txt"
    with filename as csv_file:
        msg_data = list(csv.reader(csv_file, delimiter=','))
        count = 0
        for row in msg_data:
            for i in range(0,len(row)):
                row[i] = int(row[i])

    msg_clean = [x[:] for x in [[0] * len(msg_data)] * 2]

    for i in range(0,2):
        for j in range(0,len(msg_data)):
            msg_clean[i][j] = msg_data[j][i]
            
    return msg_clean[1]

def getSameDiffClusterCount(basename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-adhoc2.txt")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    #filename = basename+"-adhoc2.txt"
    with filename as f:
        adhoc_raw_data = [line[:-1] for line in f]
    iters = getIters(basename, "-clusters.txt")
    headerCount = [0 for i in range(0,iters)]   
    index = 0

    for i in range(0,len(adhoc_raw_data)):
        if(adhoc_raw_data[i][0] is 'A'):
            try:
                headerCount[index] = int((adhoc_raw_data[i].split(" "))[1])
            except IndexError:
                #print(index)
                pass
            index = index + 1

    clusters = [[] for i in range(0,iters)]
    for j in range(0,len(clusters)):
        clusters[j] = [Cluster("") for x in range(0,headerCount[j])]
        
    index = 0
    headPerIterCount = 0
    errString = ""

    for i in range(0,len(adhoc_raw_data)):
        if(adhoc_raw_data[i][0] is not 'A'):
            line = adhoc_raw_data[i].split(":")
            line_str = "H"+line[0]+"-M"+line[1][1:]
            try:
                clusters[index][headPerIterCount] = Cluster(line_str)
                clusters[index][headPerIterCount].process()
                headPerIterCount = headPerIterCount + 1
            except IndexError:
                donothing = 1
        else:
            index = index + 1
            headPerIterCount = 0
    
    sameClusterCount = [0 for i in range(0,len(clusters)+1)]
    diffClusterCount = [0 for i in range(0,len(clusters)+1)]

    for i in range(3,len(clusters)):
        for j in range(0,len(clusters[i-1])):
            for k in range(0,len(clusters[i])):
                if(clusters[i-1][j].ch == clusters[i][k].ch):
                    s1 = clusters[i-1][j].cm
                    s2 = clusters[i][k].cm
                    if(len(set(s1).intersection(set(s2))) != 0):
                        sameClusterCount[i] = sameClusterCount[i]+1

        diffClusterCount[i] = len(clusters[i])-sameClusterCount[i] 
    
    return sameClusterCount[1:], diffClusterCount[1:]

def get4GandBTReadingCounts(basename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-cluster_readings.txt")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    #filename = basename+"-cluster_readings.txt"
    with filename as f:
        creadings = f.readlines()

    bt_count = 0
    svr_count = 0
    bt_bytes = 0
    svr_bytes = 0
    for i in range(0,len(creadings)):
        creadings[i] = creadings[i].split('-')
        creadings[i][6] = (creadings[i][6].split('\n'))[0]
        if(creadings[i][3]=='BT'):
            bt_count = bt_count + 1
            bt_bytes = bt_bytes + int(creadings[i][6])
        else:
            if(creadings[i][3]=='Server'):
                svr_count = svr_count + 1
                svr_bytes = svr_bytes + int(creadings[i][6])

    #print(bt_count)
    bt_cl_data = (bt_count,bt_bytes)
    svr_cl_data = (svr_count,svr_bytes)
    svr_cl_bytes = svr_bytes
    bt_cl_bytes = bt_bytes
    svr_cl_count = svr_count
    bt_cl_count = bt_count




    return(svr_cl_count, bt_cl_count)

def getAliveValidPercentPerIter(basename):
    zf = ZipFile(basename, 'r')
    prefix = os.path.split(basename)[1].split(".zip")[0]
    filename = "%s%s" % (prefix, "-nodes_alive_valid.txt")
    f = zf.open(filename)
    filename = StringIO(f.read().decode("utf-8")) #don't forget this line!
    #filename = basename+"-nodes_alive_valid.txt"
    with filename as f:
        raw_av_data = f.readlines()
        node_av_data = [x[:] for x in [["0"] * 2] * len(raw_av_data)]

    for i in range(0,len(raw_av_data)):
        node_av_data[i] = raw_av_data[i].split(',')
        node_av_data[i][1] = (node_av_data[i][1].split("\n"))[0]
        node_av_data[i][0] = float((node_av_data[i][0].split(":"))[1])
        node_av_data[i][1] = float((node_av_data[i][1].split(":"))[1])
        
    iters = len(raw_av_data)
    percent_alive = [0.0]*iters
    for i in range(0,iters):
        try:
            percent_alive[i] = (node_av_data[i][1]/node_av_data[i][0])*100
        except ZeroDivisionError:
            percent_alive[i] = 1
    return percent_alive