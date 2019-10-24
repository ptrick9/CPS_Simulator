import multiprocessing
import itertools
import os


def runner(queue):
    while True:
        job = queue.get()
        print("%d\\%d" % (job[1], job[2]))
        #print(job)
        command = "./simulator/Simulation "+' '.join(job[0])
        print(command)
        #os.system(command + " 1>/dev/null")
        queue.task_done()


if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()


    #switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvSensor=false -csvMove=true -zipFiles=true"]
    switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvMove=true -zipFiles=true -windRegionPath=hull_fine_bomb_fixWind_9x9.txt"]
    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt -iterations=10000 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_4.txt', 'FireflyWalls.png']]]

    paths = []
    for pop in [[200, 10000], [500, 4000], [1000, 2000], [2000, 1000], [5000, 400], [10000, 200]]:
        for it in [1, 2]:
            paths.append((pop[0], it, pop[1]))

    #movementPath = ["-movementPath=%s" % s for s in ["/home/simulator/git-simulator/movement/marathon_street_2000_%d.scb" % i for i in range(1,10)]]
    movementPath = ["-movementPath=/home/simulator/git-simulator/movement/Firefly_%d_%d.scb -totalNodes=%d -moveSize=%d" % (s[0], s[1], s[0], s[2]) for s in paths]

    sensorPath = ["-sensorPath=%s" %s for s in ["smooth_marathon.csv"]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5]]
    sittingStopThreshold = ["-sittingStopThreshold=%d" % d for d in [5]]
    negativeSittingStopThreshold = ["-negativeSittingStopThreshold=%d" % d for d in [-10]]
    GridCapacityPercentage = ["-GridCapacityPercentage=%f" % f for f in [0.9]]
    naturalLoss = ["-naturalLoss=%f" % f for f in [0.005]]
    sensorSamplingLoss = ["-sensorSamplingLoss=%f" % f for f in [0.001]]
    GPSSamplingLoss = ["-GPSSamplingLoss=%f" % f for f in [0.005]]
    serverSamplingLoss = ["-serverSamplingLoss=%f" % f for f in [0.01]]
    SamplingLossBTCM = ["-SamplingLossBTCM=%f" % f for f in [0.0001]]
    SamplingLossWifiCM = ["-SamplingLossWifiCM=%f" % f for f in [0.001]]
    SamplingLoss4GCM = ["-SamplingLoss4GCM=%f" % f for f in [0.005]]
    SamplingLossAccelCM = ["-SamplingLossAccelCM=%f" % f for f in [0.001]]
    thresholdBatteryToHave = ["-thresholdBatteryToHave=%d" % d for d in [30]]
    thresholdBatteryToUse = ["-thresholdBatteryToUse=%d" % d for d in [10]]
    movementSamplingSpeed = ["-movementSamplingSpeed=%d" % d for d in [20]]
    movementSamplingPeriod = ["-movementSamplingPeriod=%d" % d for d in [20]]
    maxBufferCapacity = ["-maxBufferCapacity=%d" % d for d in [25]]
    sensorSamplingPeriod = ["-sensorSamplingPeriod=%d" % d for d in [1000]]
    GPSSamplingPeriod = ["-GPSSamplingPeriod=%d" % d for d in [1000]]
    serverSamplingPeriod = ["-serverSamplingPeriod=%d" % d for d in [1000]]
    nodeStoredSamples = ["-nodeStoredSamples=%d" % d for d in [10]]
    gridStoredSample = ["-GridStoredSamples=%d" % d for d in [10]]
    errorMultiplier = ["-errorMultiplier=%f" % f for f in [1.76]]
    numSuperNodes = ["-numSuperNodes=%d" %d for d in [4]]
    recalibThresh = ["-RecalibrationThreshold=%d" % d for d in [3]]
    StandardDeviationThreshold = ["-StandardDeviationThreshold=%f" % f for f in [1.7]]
    SuperNodeSpeed = ["-SuperNodeSpeed=%d" % d for d in [3]]
    #Squares = ["-SquareRowCM=%d -SquareColCM=%d" % (d, d) for d in [3, 4, 5]]
    #validationThreshold = ["-validationThreshold=%d" % d for d in [0]]
    serverRecal = ["-serverRecal=%s" % s for s in ['true', 'false']]
    driftExplore = ["-driftExplorer=%s" % s for s in ['true']]
    helper = ["-fineSensorPath=%s -csvSensor=%s" % (s[0], s[1]) for s in [("fine_bomb9x9.csv", 'false')]]
    commBomb = ["-commandBomb=%s -bombX=%d -bombY=%d" % (s[0], s[1], s[2]) for s in [("true", 1007, 1491),("true", 988, 1333),("true", 213, 256),("true", 1090, 154),("true", 289, 366),("true", 1275, 875),("true", 840, 1436),("true", 808, 1149),("true", 874, 1365),("true", 983, 891),("true", 1117, 1813),("true", 1172, 1355),("true", 223, 263),("true", 1137, 923),("true", 1109, 1086),("true", 921, 1212),("true", 800, 419),("true", 976, 1170),("true", 1000, 1179),("true", 713, 383),("true", 1199, 937),("true", 266, 646),("true", 1340, 809),("true", 1147, 847),("true", 1127, 959)]]
    #detectionWindow = ["-detectionWindow=%d" % d for d in [100]]

    validTypes = ["-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=100", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=500", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=100", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=500", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=100", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=500", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=100", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=500", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=100", "-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=500", "-validationType=square -SquareRowCM=3 -SquareColCM=3", "-validationType=square -SquareRowCM=4 -SquareColCM=4", "-validationType=square -SquareRowCM=5 -SquareColCM=5"]

    recalReject = ["-recalReject=%s" % s for s in ['true', 'false']]

    #sensorPath, fineSensorPath
    runs = (list(itertools.product(*[switches, scenarios, movementPath, detectionThreshold,
                                     detectionDistance, sittingStopThreshold, negativeSittingStopThreshold, GridCapacityPercentage,
                                     naturalLoss,sensorSamplingLoss,GPSSamplingLoss,serverSamplingLoss,SamplingLossBTCM,SamplingLossWifiCM,
                                     SamplingLoss4GCM,SamplingLossAccelCM,thresholdBatteryToHave,thresholdBatteryToUse,movementSamplingSpeed,
                                     movementSamplingPeriod,maxBufferCapacity,sensorSamplingPeriod,GPSSamplingPeriod,serverSamplingPeriod,
                                     nodeStoredSamples,gridStoredSample,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,
                                     SuperNodeSpeed,serverRecal, driftExplore, commBomb, helper, validTypes, recalReject])))
    print(len(runs))

    factor = 1
    x = 0
    for r in runs:
        for i in range(factor):
            j = [zz for zz in r]
            j.append("-OutputFileName=/home/simulator/simData/driftExplorerFireflyBombFinal/Log_%d" % x)

            
            v = j
            q.put([v, x, len(runs)*factor])
            x+= 1

       
    p = multiprocessing.Pool(20, runner, (q,))

    q.join()
