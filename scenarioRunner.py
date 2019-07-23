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
        #print(job)
        command = "./simulator "+' '.join(job)
        print(command)
        #os.system(command)
        queue.task_done()

if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()


    switches = ["-logNodes=false -logPostion=true -logGrid=false -logEnergy=false -regionRouting=true -clusteringOn=true -noEnergy=true -csvSensor=true -csvMove=true"]

    scenarios = ["-inputFileName=%s -imageFileName=%s -stimFileName=circle_0.txt -outRoutingStatsName=routingStats.txt   \
                -iterations=1000 -superNodes=false -doOptimize=false" % (s[0], s[1]) for s in [['Scenario_3.txt', 'marathon_street_map.png']]]

    row = ["-squareRow=%d" % d for d in [60, 120]]
    col = ["-squareCXol=%d" % d for d in [320, 640]]
    movementPath = ["-movementPath=%s" % s for s in ["marathon_street_2k.txt"]]
    sensorPath = ["-sensorPath=%s" %s for s in ["smoothed_marathon.csv"]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5, 6, 7]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6, 7]]
    outputFileName = ["-p.OutputFileName=%s" % s for s in ["Log"]]
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
    gridStoredSample = ["-p.GridStoredSample=%d" % d for d in [10]]
    errorMultiplier = ["-errorMultiplier=%f" % f for f in [1.0]]
    numSuperNodes = ["-numSuperNodes=%d" %d for d in [4]]
    recalibThresh = ["-Recalibration Threshold=%d" % d for d in [3]]
    StandardDeviationThreshold = ["-StandardDeviationThreshold=%f" % f for f in [1.7]]
    SuperNodeSpeed = ["-SuperNodeSpeed=%d" % d for d in [3]]
    SquareRowCM = ["-SquareRowCM=%d" % d for d in [60]]
    SquareColCM = ["-SquareColCM=%d" % d for d in [320]]
    validationThreshold = ["-validationThreshold=%d" % d for d in [1]]


    runs = (list(itertools.product(*[scenarios, row, col, movementPath, detectionThreshold, detectionDistance, outputFileName, sittingStopThreshold, negativeSittingStopThreshold, GridCapacityPercentage, naturalLoss,sensorSamplingLoss,GPSSamplingLoss,serverSamplingLoss,SamplingLossBTCM,SamplingLossWifiCM,SamplingLoss4GCM,SamplingLossAccelCM,thresholdBatteryToHave,thresholdBatteryToUse,movementSamplingSpeed,movementSamplingPeriod,maxBufferCapacity,sensorSamplingPeriod,GPSSamplingPeriod,serverSamplingPeriod,nodeStoredSamples,gridStoredSample,detectionThreshold,errorMultiplier,numSuperNodes,recalibThresh,StandardDeviationThreshold,detectionDistance,SuperNodeSpeed,SquareRowCM,SquareColCM,validationThreshold])))
    
    x = 0
    for r in runs:
        for i in range(10):
            j = [zz for zz in r]
            j.append("-outputFileName=big_data/Log_%d" % x)
            v = j
            q.put(v)
            x+= 1

       
    p = multiprocessing.Pool(40, runner, (q,))

    q.join()
