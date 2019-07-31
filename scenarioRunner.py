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


    switches = ["-logNodes=false -logPosition=true -logGrid=false -logEnergy=false -regionRouting=true -noEnergy=true -csvSensor=true -csvMove=true"]

    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt -iterations=5000 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_3.txt', 'marathon_street_map.png']]]


    movementPath = ["-movementPath=%s" % s for s in ["/home/simulator/git-simulator/movement/marathon_street_2000_%d.scb" % i for i in range(1,11)]]
    sensorPath = ["-sensorPath=%s" %s for s in ["smooth_marathon.csv"]]
    fineSensorPath = ["-fineSensorPath=%s" %s for s in ["fine_bomb_marathon.csv"]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6]]
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
    errorMultiplier = ["-errorMultiplier=%f" % f for f in [1.0]]
    numSuperNodes = ["-numSuperNodes=%d" %d for d in [4]]
    recalibThresh = ["-RecalibrationThreshold=%d" % d for d in [3]]
    StandardDeviationThreshold = ["-StandardDeviationThreshold=%f" % f for f in [1.7]]
    SuperNodeSpeed = ["-SuperNodeSpeed=%d" % d for d in [3]]
    SquareRowCM = ["-SquareRowCM=%d" % d for d in [60]]
    SquareColCM = ["-SquareColCM=%d" % d for d in [320]]
    validationThreshold = ["-validationThreshold=%d" % d for d in [0, 1, 2, 3, 4, 5]]


    runs = (list(itertools.product(*[switches, scenarios, movementPath, sensorPath, fineSensorPath, detectionThreshold, detectionDistance, sittingStopThreshold, negativeSittingStopThreshold, GridCapacityPercentage, naturalLoss,sensorSamplingLoss,GPSSamplingLoss,serverSamplingLoss,SamplingLossBTCM,SamplingLossWifiCM,SamplingLoss4GCM,SamplingLossAccelCM,thresholdBatteryToHave,thresholdBatteryToUse,movementSamplingSpeed,movementSamplingPeriod,maxBufferCapacity,sensorSamplingPeriod,GPSSamplingPeriod,serverSamplingPeriod,nodeStoredSamples,gridStoredSample,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,SuperNodeSpeed,SquareRowCM,SquareColCM,validationThreshold])))
    
    x = 0
    for r in runs:
        for i in range(10):
            j = [zz for zz in r]
            j.append("-OutputFileName=/home/simulator/simData/fineGrainedBomb/Log_%d" % x)
            v = j
            q.put([v, x, len(runs)*10])
            x+= 1

       
    p = multiprocessing.Pool(7, runner, (q,))

    q.join()
