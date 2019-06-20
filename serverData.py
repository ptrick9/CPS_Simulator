import numpy as np
#from numpy import *
import pylab as pl
import re
import pandas as pd

def main():
    stats = []
    meanNum = []
    stdDevNum = []
    varianceNum = []

    nodeNums = []

    points = []
    calTimes = []
    sensitivity = []
    noise = []
    error = []
    
    log = open("Log-server.txt", 'r')
    if log.mode == 'r':
        stats = log.readlines()
    log.close()
        
    stats[1] = re.sub("[\[\\n\]]", "", stats[1])
    stats[3] = re.sub("[\[\\n\]]", "", stats[3])
    stats[5] = re.sub("[\[\\n\]]", "", stats[5])
    
    mean = re.split(" ", stats[1])
    mean = mean[1:]
    #print(mean)
    for i in mean:
        meanNum.append(float(i))

    stdDev = re.split(" ", stats[3])
    stdDev = stdDev[1:]
    #print(stdDev)
    for i in stdDev:
        stdDevNum.append(float(i))

    variance = re.split(" ", stats[5])
    variance = variance[1:]
    #print(variance)
    for i in variance:
        varianceNum.append(float(i))

    t = np.linspace(0, 70, num = 70)
    #t2 = np.linspace(0, 1000, num = 1000)

    vals = open("Log-nodeTest.txt", 'r')
    if vals.mode == 'r':
        line = vals.readlines()
    vals.close()
   
    for i in range(len(line)):
        if i % 4 == 0:
            nodeNums.append(float(re.sub("[Val: \\n]", "",line[i])))
        if i % 4 == 1:
            sensitivity.append(float(re.sub("[Sensi: \\n]", "",line[i])))
        if i % 4 == 2:
            noise.append(float(re.sub("[Noise: \\n]", "",line[i])))
        if i % 4 == 3:
            error.append(float(re.sub("[Error: \\n]", "",line[i])))

    t2 = np.linspace(0, len(nodeNums), num = len(nodeNums))
    window = 8

    df = pd.DataFrame(nodeNums)
    rollingMean = df.rolling(window).mean()

    df = pd.DataFrame(rollingMean)
    rollingMin = df.rolling(window).min()
    rollingMax = df.rolling(window).max()

    calib = open("Log-nodeTest2.txt", 'r')
    if calib.mode == 'r':
       calVal  = calib.readlines()
    calib.close()
    
    '''calVal[0] = re.sub("[\[\]\\n]", "", calVal[0])
    calVal[0] = calVal[0][1:]
    times = re.split(" ", calVal[0])
    for i in times:
        calTimes.append(float(i))
    print(calTimes)
    
    calVal[1] = re.sub("[\[\]\\n]", "", calVal[1])
    calVal[1] = calVal[1][1:]
    calReads = re.split(" ", calVal[1])
    for i in calReads:
        points.append(float(i))'''

    pl.plt.title("Mean")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , meanNum, "r")
    pl.plt.show()
    
    pl.plt.title("Variance")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , varianceNum, "g")
    pl.plt.show()

    pl.plt.title("Standard Deviation")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , stdDevNum, "b")
    pl.plt.show()
    
    pl.plt.title("Readings")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , nodeNums, "r", label='Sensor Values')
    #pl.plt.plot(t2 , noise, "b", label='Noise')
    #pl.plt.legend()
    #pl.plt.show()

    pl.plt.title("Readings")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , rollingMean, "b", label='Rolling Mean')
    pl.plt.legend()
    pl.plt.show()

    pl.plt.title("Min and Max")
    pl.plt.xlabel("Time (Iterations)")
    pl.plt.plot(t2 , rollingMin, "r", label='Min')
    pl.plt.plot(t2 , rollingMax, "b", label='Max')
    pl.plt.legend()
    pl.plt.show()

    pl.plt.plot(t2 , rollingMax - rollingMin, "b", label='Max')
    pl.plt.show()
    
    pl.plt.title("Sensitivity and Error")
    pl.plt.plot(t2 , sensitivity, "g", label='Sensitivity')
    pl.plt.plot(t2 , error, "y", label='Error')
    pl.plt.legend()
    pl.plt.show()

if __name__== "__main__":
    main()
