import multiprocessing
import itertools
import os

'''
-inputFileName=Scenario_3.txt
-imageFileName=marathon_street_map.png
-logPosition=true
-logGrid=true
-logEnergy=true
-logNodes=false
-noEnergy=true
-sensorPath=C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/smoothed_marathon.csv
-SquareRowCM=60
-SquareColCM=320
-csvMove=true
-movementPath=marathon_street_2k.txt
-iterations=1000
-csvSensor=true
-detectionThreshold=5
-superNodes=false
-detectionDistance=6
-cpuprofile=event

'''

def runner(queue):
    while True:
        job = queue.get()
        print("%d\\%d" % (job[1], job[2]))
        #print(job)
        command = "./simulator/Simulation "+' '.join(job[0])
        #print(command)
        os.system(command + " 1>/dev/null")
        queue.task_done()


if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()


    #switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvSensor=false -csvMove=true -zipFiles=true"]
    switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvMove=true -zipFiles=true -windRegionPath=hull_fine_bomb_9x9.txt"]
    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt -iterations=10000 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_3.txt', 'marathon_street_map.png']]]

    paths = []
    for pop in [[200, 10000], [500, 4000], [1000, 2000], [2000, 1000], [5000, 400], [10000, 200]]:
        for it in [1, 2, 3, 4]:
            paths.append((pop[0], it, pop[1]))

    #movementPath = ["-movementPath=%s" % s for s in ["/home/simulator/git-simulator/movement/marathon_street_2000_%d.scb" % i for i in range(1,10)]]
    movementPath = ["-movementPath=/home/simulator/git-simulator/movement/marathon2_%d_%d.scb -totalNodes=%d -moveSize=%d" % (s[0], s[1], s[0], s[2]) for s in paths]

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
    SquareRowCM = ["-SquareRowCM=%d" % d for d in [60]]
    SquareColCM = ["-SquareColCM=%d" % d for d in [320]]
    validationThreshold = ["-validationThreshold=%d" % d for d in [0, 1, 2, 3, 4]]
    serverRecal = ["-serverRecal=%s" % s for s in ['true', 'false']]
    driftExplore = ["-driftExplorer=%s" % s for s in ['false']]
    helper = ["-fineSensorPath=%s -csvSensor=%s" % (s[0], s[1]) for s in [("fine_bomb9x9.csv", 'false')]]
    commBomb = ["-commandBomb=%s -bombX=%d -bombY=%d" % (s[0], s[1], s[2]) for s in [('true', 1197, 132),('true', 339, 213),('true', 1215, 210),('true', 1115, 210),('true', 708, 134),('true', 1048, 146),('true', 318, 237),('true', 968, 214),('true', 1491, 147),('true', 1171, 132),('true', 577, 213),('true', 197, 215),('true', 582, 211),('true', 1286, 214),('true', 545, 234),('true', 947, 211),('true', 281, 222),('true', 1445, 149),('true', 997, 213),('true', 1483, 213)]]
    detectionWindow = ["-detectionWindow=%d" % d for d in [100, 500, 1000]]

    #sensorPath, fineSensorPath
    runs = (list(itertools.product(*[switches, scenarios, movementPath, detectionThreshold,
                                     detectionDistance, sittingStopThreshold, negativeSittingStopThreshold, GridCapacityPercentage,
                                     naturalLoss,sensorSamplingLoss,GPSSamplingLoss,serverSamplingLoss,SamplingLossBTCM,SamplingLossWifiCM,
                                     SamplingLoss4GCM,SamplingLossAccelCM,thresholdBatteryToHave,thresholdBatteryToUse,movementSamplingSpeed,
                                     movementSamplingPeriod,maxBufferCapacity,sensorSamplingPeriod,GPSSamplingPeriod,serverSamplingPeriod,
                                     nodeStoredSamples,gridStoredSample,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,
                                     SuperNodeSpeed,SquareRowCM,SquareColCM,validationThreshold,serverRecal, driftExplore, commBomb, helper, detectionWindow])))
    factor = 1
    x = 0
    for r in runs:
        for i in range(factor):
            j = [zz for zz in r]
            j.append("-OutputFileName=/home/simulator/simData/driftExplorerBombADC_more/Log_%d" % x)

            
            v = j
            q.put([v, x, len(runs)*factor])
            x+= 1

       
    p = multiprocessing.Pool(20, runner, (q,))

    q.join()
