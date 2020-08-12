import multiprocessing
import itertools
import os


def runner(queue):
    while True:
        job = queue.get()
        print("%d\\%d" % (job[1], job[2]))
        #print(job)
        command = "./simulator/simulator "+' '.join(job[0])
        #print(command)
        os.system(command + " 1>/dev/null")
        queue.task_done()


if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()


    #switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvSensor=false -csvMove=true -zipFiles=true"]
    switches = ["-logNodes=false -logPosition=false -logGrid=false -logOutput=false -logEnergy=false -logClusters=true -logBattery=true -regionRouting=true -csvMove=true -zipFiles=true -windRegionPath=hull_fine_bomb_fixWind_9x9.txt"]
    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt -iterations=4500 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_Stadium.txt', 'DelawareStadiumWalls.png']]]#,['Scenario_3.txt', 'marathon_street_map.png']]]

    paths = []
    for pop in [[500, 10000], [1000, 1000], [2000, 1000], [3500, 500]]: #[[200, 10000], [500, 4000], [1000, 2000], [2000, 1000], [5000, 400], [10000, 200]]:
        for it in [1, 2, 3, 4]:
            paths.append((pop[0], it, pop[1]))

    #movementPath = ["-movementPath=%s" % s for s in ["/home/simulator/git-simulator/movement/marathon_street_2000_%d.scb" % i for i in range(1,10)]]
    #movementPath = ["-movementPath=/home/simulator/git-simulator/movement/marathon2_%d_%d.scb -totalNodes=%d -moveSize=%d" % (s[0], s[1], s[0], s[2]) for s in paths]
    movementPath = ["-movementPath=/home/simulator/git-simulator/movement/Stadium_%d_%d.scb -totalNodes=%d -moveSize=%d" % (s[0], s[1], s[0], s[2]) for s in paths]

    sensorPath = ["-sensorPath=%s" %s for s in ["smooth_marathon.csv"]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5]]
    sittingStopThreshold = ["-sittingStopThreshold=%d" % d for d in [5]]
    #maxBufferCapacity = ["-maxBufferCapacity=%d" % d for d in [25]]
    nodeStoredSamples = ["-nodeStoredSamples=%d" % d for d in [10]]
    errorMultiplier = ["-errorMultiplier=%f" % f for f in [1.76]]
    numSuperNodes = ["-numSuperNodes=%d" %d for d in [4]]
    recalibThresh = ["-RecalibrationThreshold=%d" % d for d in [3]]
    StandardDeviationThreshold = ["-StandardDeviationThreshold=%f" % f for f in [1.7]]
    superNodes = ["-superNodes=%s" % s for s in ['false']]
    #Squares = ["-SquareRowCM=%d -SquareColCM=%d" % (d, d) for d in [3, 4, 5]]
    #validationThreshold = ["-validationThreshold=%d" % d for d in [0]]
    serverRecal = ["-serverRecal=%s" % s for s in ['true']]
    driftExplore = ["-driftExplorer=%s" % s for s in ['false']]
    helper = ["-fineSensorPath=%s -csvSensor=%s" % (s[0], s[1]) for s in [("fine_bomb9x9.csv", 'false')]]
    #commBomb = ["-commandBomb=%s -bombX=%d -bombY=%d" % (s[0], s[1], s[2]) for s in [("true", 848, 145),("true", 191, 139),("true", 1570, 143),("true", 256, 142),("true", 211, 128),("true", 897, 216),("true", 1290, 144),("true", 985, 126),("true", 456, 128),("true", 813, 138),("true", 743, 212),("true", 1177, 148),("true", 379, 236),("true", 1126, 225),("true", 729, 146),("true", 909, 141),("true", 558, 149),("true", 825, 144),("true", 439, 217),("true", 556, 219),("true", 919, 130),("true", 673, 211),("true", 668, 214),("true", 651, 210),("true", 907, 124)]]
    commBomb = ["-commandBomb=%s" % s for s in ['false']]

    #detectionWindow = ["-detectionWindow=%d" % d for d in [100]]

    #validTypes = ['-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=50', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=100', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=200', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=50', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=100', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=200', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=50', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=100', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=200', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=50', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=100', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=200', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=50', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=100', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=200', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=5', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=10', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=20', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=5', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=10', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=20', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=5', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=10', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=20', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=5', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=10', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=20', '-validationType=square -SquareRowCM=10 -SquareColCM=10 -GridStoredSamples=5', '-validationType=square -SquareRowCM=10 -SquareColCM=10 -GridStoredSamples=10', '-validationType=square -SquareRowCM=10 -SquareColCM=10 -GridStoredSamples=20', '-validationType=square -SquareRowCM=12 -SquareColCM=12 -GridStoredSamples=5', '-validationType=square -SquareRowCM=12 -SquareColCM=12 -GridStoredSamples=10', '-validationType=square -SquareRowCM=12 -SquareColCM=12 -GridStoredSamples=20', '-validationType=square -SquareRowCM=14 -SquareColCM=14 -GridStoredSamples=5', '-validationType=square -SquareRowCM=14 -SquareColCM=14 -GridStoredSamples=10', '-validationType=square -SquareRowCM=14 -SquareColCM=14 -GridStoredSamples=20', '-validationType=square -SquareRowCM=16 -SquareColCM=16 -GridStoredSamples=5', '-validationType=square -SquareRowCM=16 -SquareColCM=16 -GridStoredSamples=10', '-validationType=square -SquareRowCM=16 -SquareColCM=16 -GridStoredSamples=20']
    #validTypes = ['-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=60', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=120', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=0 -detectionWindow=240', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=60', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=120', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=1 -detectionWindow=240', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=60', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=120', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=2 -detectionWindow=240', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=60', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=120', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=3 -detectionWindow=240', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=60', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=120', '-validationType=validators -SquareRowCM=3 -SquareColCM=3 -validationThreshold=4 -detectionWindow=240', '-validationType=square -SquareRowCM=1 -SquareColCM=1 -GridStoredSamples=5', '-validationType=square -SquareRowCM=1 -SquareColCM=1 -GridStoredSamples=10', '-validationType=square -SquareRowCM=1 -SquareColCM=1 -GridStoredSamples=20', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=5', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=10', '-validationType=square -SquareRowCM=2 -SquareColCM=2 -GridStoredSamples=20', '-validationType=square -SquareRowCM=3 -SquareColCM=3 -GridStoredSamples=5', '-validationType=square -SquareRowCM=3 -SquareColCM=3 -GridStoredSamples=10', '-validationType=square -SquareRowCM=3 -SquareColCM=3 -GridStoredSamples=20', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=5', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=10', '-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=20', '-validationType=square -SquareRowCM=5 -SquareColCM=5 -GridStoredSamples=5', '-validationType=square -SquareRowCM=5 -SquareColCM=5 -GridStoredSamples=10', '-validationType=square -SquareRowCM=5 -SquareColCM=5 -GridStoredSamples=20', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=5', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=10', '-validationType=square -SquareRowCM=6 -SquareColCM=6 -GridStoredSamples=20', '-validationType=square -SquareRowCM=7 -SquareColCM=7 -GridStoredSamples=5', '-validationType=square -SquareRowCM=7 -SquareColCM=7 -GridStoredSamples=10', '-validationType=square -SquareRowCM=7 -SquareColCM=7 -GridStoredSamples=20', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=5', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=10', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=20']
    validTypes = ['-validationType=square']
    recalReject = ["-recalReject=%s" % s for s in ['false']]

    clusteringOn = ["-clusteringOn=%s" % s for s in ['true']]
    clusterMax = ["-clusterMaxThresh=%d" % d for d in [40]]
    clusterMin = ["-clusterMinThresh=%d" % d for d in [0]]
    nodeBTRange = ["-nodeBTRange=%d" % d for d in [20]]
    degreeWeight = ["-degreeWeight=%f" % f for f in [0.6]]
    batteryWeight = ["-batteryWeight=%f" % f for f in [0.4]]
    aloneRCThreshold = ["-aloneThreshold=%f -aloneOrClusterRatio=true" % f for f in [0.04]]
    clusterRCThreshold = ["-aloneThreshold=%f -aloneOrClusterRatio=false" % f for f in []]
    aloneThreshold = aloneRCThreshold + clusterRCThreshold
    aloneClusterSearchOptions = (list(itertools.product(*[['-aloneClusterSearch=true'], clusterMin])))
    for i in range(len(aloneClusterSearchOptions)):
        aloneClusterSearchOptions[i] = ' '.join(aloneClusterSearchOptions[i])
    aloneClusterSearchOptions += ['-aloneClusterSearch=false']
    reclusterPeriod = ["-reclusterPeriod=%d" % d for d in [200]]
    localRecluster = ["-localRecluster=%d -expansiveRatio=%f" % (d, f) for (d, f) in [(0,0), (2,0)]]
    batteryCap = ["-batteryCapacity=%d" % d for d in [100000000]]
    losses = ["-bluetoothLossPercentage=%f -wifiLossPercentage=%f -sampleLossPercentage=%f" % (f[0],f[1],f[2]) for f in [(0.00001,0.0001,0.00002)]]
    clusterSearchOptions = ["-clusterSearchThresh=%d -adaptiveClusterSearch=%s" % (d,s) for (d,s) in [(4, 'true')]]
    #initClusterTime = ["-initClusterTime=%d" % d for d in [0, 200, 600]]
    CHTimeThresh = ["-CHTimeThresh=%d" % d for d in [450, 600, 900]]#600
    CHBatteryDropThresh = ["-CHBatteryDropThresh=%f" % f for f in [0.2,0.3,0.4]]
    smallImprovementRatio = ["-smallImprovementRatio=%f" % f for f in [0.33]]
    largeImprovement = ["-largeImprovement=%f" % f for f in [0.75]]
    GRIncrementDecrement = ["-GRIncrement=%f -GRDecrement=%f" % (f[0],f[1]) for f in [(1.4,0.75)]]
    disableGRThresh = ["-disableGRThresh=%f" % f for f in [0.9]]
    disableCSThresh = ["-disableCSThresh=%f" % f for f in [0]]
    maxClusterHeads = ["-maxClusterHeads=%d" % d for d in [1]]
    serverReadyThresh = ["-serverReadyThresh=%f" % f for f in [0.98]]

    globalReclusterOnOptions = (list(itertools.product(*[smallImprovementRatio, largeImprovement, GRIncrementDecrement,
                                                        serverReadyThresh, disableGRThresh])))

    for i in range(len(globalReclusterOnOptions)):
        globalReclusterOnOptions[i] = ' '.join(globalReclusterOnOptions[i])

    globalRecluster0 = (list(itertools.product(*[["-globalRecluster=0"], aloneClusterSearchOptions])))
    for i in range(len(globalRecluster0)):
        globalRecluster0[i] = ' '.join(globalRecluster0[i])

    globalRecluster1 = (list(itertools.product(*[["-globalRecluster=1 -aloneClusterSearch=false"], aloneThreshold, globalReclusterOnOptions])))
    for i in range(len(globalRecluster1)):
        globalRecluster1[i] = ' '.join(globalRecluster1[i])

    globalRecluster2 = (list(itertools.product(*[["-globalRecluster=2 -aloneClusterSearch=false"], reclusterPeriod, globalReclusterOnOptions])))
    for i in range(len(globalRecluster2)):
        globalRecluster2[i] = ' '.join(globalRecluster2[i])

    globalRecluster3 = (list(itertools.product(*[["-globalRecluster=3 -aloneClusterSearch=false"], aloneThreshold, reclusterPeriod, globalReclusterOnOptions])))
    for i in range(len(globalRecluster3)):
        globalRecluster3[i] = ' '.join(globalRecluster3[i])

    globalRecluster4 = (list(itertools.product(*[["-globalRecluster=4 -aloneClusterSearch=false"], aloneThreshold, reclusterPeriod, globalReclusterOnOptions])))
    for i in range(len(globalRecluster4)):
        globalRecluster4[i] = ' '.join(globalRecluster4[i])

    globalReclusterOptions = globalRecluster0 + globalRecluster1 + globalRecluster2

    #globalReclusterOptions += ["-globalRecluster=0 -aloneThreshold=0.05 -reclusterPeriod=300.0"]
    #globalReclusterOptions += ["-globalRecluster=2 -aloneThreshold=0.05 -reclusterPeriod=25.0 -GRIncrement=1 -GRDecrement=1 -disableGRThresh=0.7 -serverReadyThresh=0.995"]
    #globalReclusterOptions += ["-globalRecluster=-1 -aloneThreshold=0.05 -reclusterPeriod=300.0"]

    clusteringOptions = (list(itertools.product(*[clusteringOn, clusterMax, clusterMin, degreeWeight, batteryWeight,
                                                globalReclusterOptions, clusterSearchOptions, maxClusterHeads,
                                                disableCSThresh, localRecluster, CHTimeThresh, CHBatteryDropThresh])))
    for i in range(len(clusteringOptions)):
        clusteringOptions[i] = ' '.join(clusteringOptions[i])
    clusteringOptions += ["-clusteringOn=false"]

    #sensorPath, fineSensorPath
    runs = (list(itertools.product(*[switches, scenarios, movementPath, detectionThreshold,
                                    detectionDistance, sittingStopThreshold, batteryCap,
                                    nodeStoredSamples,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,
                                    superNodes,serverRecal, driftExplore, commBomb, helper, validTypes, recalReject,
                                    nodeBTRange, losses, clusteringOptions])))
    print(len(runs))

    factor = 1
    x = 0
    for r in runs:
        for i in range(factor):
            j = [zz for zz in r]
            j.append("-OutputFileName=/home/simulator/git-simulator/CPS_Simulator/simData/ReportSmallTests4/Log_%d" % x)
            
            v = j
            q.put([v, x, len(runs)*factor])
            x+= 1

       
    p = multiprocessing.Pool(23, runner, (q,))

    q.join()
