import multiprocessing
import itertools
import os


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
    switches = ["-logNodes=false -logPosition=false -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvMove=true -zipFiles=true -windRegionPath=hull_fine_bomb_fixWind_9x9.txt"]
    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt -iterations=10000 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_3.txt', 'marathon_street_map.png']]]

    paths = []
    for pop in [[500, 4000], [2000, 1000], [10000, 200]]:
        for it in [3, 4]:
            paths.append((pop[0], it, pop[1]))

    #movementPath = ["-movementPath=%s" % s for s in ["/home/simulator/git-simulator/movement/marathon_street_2000_%d.scb" % i for i in range(1,10)]]
    movementPath = ["-movementPath=/home/simulator/git-simulator/movement/marathon2_%d_%d.scb -totalNodes=%d -moveSize=%d" % (s[0], s[1], s[0], s[2]) for s in paths]

    sensorPath = ["-sensorPath=%s" %s for s in ["smooth_marathon.csv"]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5]]
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
    #commBomb = ["-commandBomb=%s -bombX=%d -bombY=%d" % (s[0], s[1], s[2]) for s in [("true", 848, 145),("true", 191, 139),("true", 1570, 143),("true", 256, 142),("true", 211, 128),("true", 897, 216),("true", 1290, 144),("true", 985, 126),("true", 456, 128),("true", 813, 138),("true", 743, 212),("true", 1177, 148),("true", 379, 236),("true", 1126, 225),("true", 729, 146),("true", 909, 141),("true", 558, 149),("true", 825, 144),("true", 439, 217),("true", 556, 219),("true", 919, 130),("true", 673, 211),("true", 668, 214),("true", 651, 210),("true", 907, 124)]]
    commBomb = ["-commandBomb=%s" % s for s in ['false']]

    #detectionWindow = ["-detectionWindow=%d" % d for d in [100]]

    validTypes = ['-validationType=square -SquareRowCM=4 -SquareColCM=4 -GridStoredSamples=10', '-validationType=square -SquareRowCM=8 -SquareColCM=8 -GridStoredSamples=10', '-validationType=square -SquareRowCM=12 -SquareColCM=12 -GridStoredSamples=10', '-validationType=square -SquareRowCM=16 -SquareColCM=16 -GridStoredSamples=10']

    recalReject = ["-recalReject=%s" % s for s in ['false']]
    #nodeBuffer = ["-nodeBuffer=%s -nodeBufferSamples=%s -nodeTimeBuffer=%s" % (s[0], s[1], s[2]) for s in [('true', 10, 10), ('true', 20, 20), ('true', 50, 50), ('false', 10, 10) ]]

    nodeBuffer = ['-nodeBuffer=true -nodeBufferSamples=10 -nodeTimeBuffer=10 -accelSendThresh=0.200000', '-nodeBuffer=true -nodeBufferSamples=10 -nodeTimeBuffer=10 -accelSendThresh=0.400000', '-nodeBuffer=true -nodeBufferSamples=10 -nodeTimeBuffer=10 -accelSendThresh=0.600000', '-nodeBuffer=true -nodeBufferSamples=10 -nodeTimeBuffer=10 -accelSendThresh=0.800000', '-nodeBuffer=true -nodeBufferSamples=10 -nodeTimeBuffer=10 -accelSendThresh=1.000000', '-nodeBuffer=true -nodeBufferSamples=20 -nodeTimeBuffer=20 -accelSendThresh=0.200000', '-nodeBuffer=true -nodeBufferSamples=20 -nodeTimeBuffer=20 -accelSendThresh=0.400000', '-nodeBuffer=true -nodeBufferSamples=20 -nodeTimeBuffer=20 -accelSendThresh=0.600000', '-nodeBuffer=true -nodeBufferSamples=20 -nodeTimeBuffer=20-accelSendThresh=0.800000', '-nodeBuffer=true -nodeBufferSamples=20 -nodeTimeBuffer=20 -accelSendThresh=1.000000', '-nodeBuffer=true -nodeBufferSamples=50 -nodeTimeBuffer=50 -accelSendThresh=0.200000', '-nodeBuffer=true -nodeBufferSamples=50 -nodeTimeBuffer=50 -accelSendThresh=0.400000', '-nodeBuffer=true -nodeBufferSamples=50 -nodeTimeBuffer=50 -accelSendThresh=0.600000', '-nodeBuffer=true -nodeBufferSamples=50 -nodeTimeBuffer=50 -accelSendThresh=0.800000', '-nodeBuffer=true -nodeBufferSamples=50 -nodeTimeBuffer=50 -accelSendThresh=1.000000', '-nodeBuffer=false -nodeBufferSamples=10 -nodeTimeBuffer=10']

    #sensorPath, fineSensorPath
    runs = (list(itertools.product(*[switches, scenarios, movementPath, detectionThreshold,
                                     detectionDistance, GridCapacityPercentage,
                                     naturalLoss,sensorSamplingLoss,GPSSamplingLoss,serverSamplingLoss,SamplingLossBTCM,SamplingLossWifiCM,
                                     SamplingLoss4GCM,SamplingLossAccelCM,thresholdBatteryToHave,thresholdBatteryToUse,movementSamplingSpeed,
                                     movementSamplingPeriod,maxBufferCapacity,sensorSamplingPeriod,GPSSamplingPeriod,serverSamplingPeriod,
                                     nodeStoredSamples,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,
                                     SuperNodeSpeed,serverRecal, driftExplore, commBomb, helper, validTypes, recalReject, nodeBuffer])))

    print(len(runs))

    factor = 1
    x = 0
    for r in runs:
        for i in range(factor):
            j = [zz for zz in r]
            j.append("-OutputFileName=/home/simulator/simData/driftExplorerSquareEnergyNoBombNetworkAccel/Log_%d" % x)

            
            v = j
            q.put([v, x, len(runs)*factor])
            x+= 1

       
    p = multiprocessing.Pool(20, runner, (q,))

    q.join()
