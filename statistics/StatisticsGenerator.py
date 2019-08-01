from statPackage.DetectionStats import *
from statPackage.ParamProcessing import *

import os
import itertools
import matplotlib.pyplot as plt
from zipfile import *

#basePath = "C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/data/bigData/"
#basePath = "C:/Users/patrick/Downloads/bigData/"
basePath = "C:/Users/patrick/Downloads/fineGrainedBomb/fineGrainedBomb/"
basePath = "C:/Users/patrick/Downloads/lowADC/"
figurePath = "C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2018/git_simulator/CPS_Simulator/g2/"

X_VAL = ['detectionThreshold', 'detectionDistance']
X_VAL = ['validationThreshold']

IGNORE = ['movementPath', 'bombX', 'bombY']
ZIP = True


def buildDetectionList(basePath, runs):
    runData = []
    for r in runs:
        runData.append(BuildDetections("%s%s" % (basePath, r)))

    for detections in runData:
        print(next(x for x in detections if x.TPConf()))


def getDetections(basePath, run):
    return BuildDetections("%s%s" % (basePath, run))


def determineDifferences(basePath, runs):
    params = {}
    for r in runs:
        p = getParameters("%s%s" % (basePath, r))
        #print(p['validationThreshold'])
        for k in p.keys():
            if 'file' not in k and 'File' not in k and k not in IGNORE:
                if k in params:
                    params[k].add(p[k])
                else:
                    params[k] = set()
                    params[k].add(p[k])

    unique = {}
    for k in params.keys():
        if len(params[k]) > 1:
            unique[k] = list(sorted(params[k]))
    #print(unique)
    return unique



# Generation of Data
def generateData(uniqueRuns):
    data = {}

    for i,file in enumerate(uniqueRuns):
        run = {}
        p = getParameters("%s%s" % (basePath, file))
        det = getDetections(basePath, file)
        firstDet = 5000
        try:
            firstDet = next(x for x in det if x.TPConf()).time
        except:
            pass
        run['Detection Time'] = firstDet
        run['# True Positive Rejections'] = sum([1 if x.TPRej() and x.time < firstDet else 0 for x in det])
        run['# False Positives'] = sum([1 if x.FP() and x.time < firstDet else 0 for x in det])
        run['# False Positive Rejections'] = sum([1 if x.FPRej() and x.time < firstDet else 0 for x in det])
        run['# False Negatives'] = sum([1 if x.FN() and x.time < firstDet else 0 for x in det])
        run['# Total False Negatives'] = sum([1 if x.FN() else 0 for x in det])
        run['# Total False Positives'] = sum([1 if x.FP() else 0 for x in det])

        run['config'] = {}
        for k in p.keys():
            run['config'][k] = p[k]

        k = tuple(p[x] for x in order)
        #k = p[X_VAL]
        if k in data:
            data[k].append(run)
        else:
            data[k] = [run]


        print("\r%d\%d" % (i, len(uniqueRuns)), end='')
    print('')
    print(data)
    return data


def generateGraphs(order, xx):
    grouped = {}
    for o in order:
        grouped[o] = {}
    for v in xx:
        for key, value in dict(zip(order, v)).items():
            if value not in grouped[key]:
                grouped[key][value] = [v]
            else:
                grouped[key][value].append(v)
    print('grouped')
    print(grouped)

    graphs = []

    if len(order) > 1:
        for choice in X_VAL:
            choices = []
            for key in grouped.keys():
                if choice != key:
                    setting = []
                    #print(grouped[key])
                    for k in grouped[key].keys():
                        setting.append({key: k})
                    choices.append(setting)

            #print('-------')
            #print(choices)
            validCombos = list(itertools.product(*choices))
            #print(validCombos)

            for combo in validCombos:
                sets = []
                descriptor = {}
                descriptor['~'] = choice
                for comboKey in combo:
                    #print(comboKey)
                    k = list(comboKey.keys())[0]
                    v = comboKey[k]
                    descriptor[k] = v
                    sets.append(set(grouped[k][v]))
                #print(sets)
                graphs.append([descriptor, set.intersection(*sets)])
        #print(graphs)
    else:
        descriptor = {}
        descriptor['~'] = order[0]
        graphs.append([descriptor, xx])
    return graphs



if __name__ == '__main__':
    #'''
    uniqueRuns = set()
    for file in os.listdir(basePath):
        uniqueRuns.add(file.split('-')[0])
    uniqueRuns = list(uniqueRuns)#[:200]

    shiftingParameters = determineDifferences(basePath, uniqueRuns)
    print(shiftingParameters)

    vals = []
    order = []
    for k in shiftingParameters.keys():
        order.append(k)
        vals.append(shiftingParameters[k])

    print(order)
    xx = list(itertools.product(*vals))
    print(xx)

    data = generateData(uniqueRuns)
    graphs = generateGraphs(order, xx)
    #'''
    #print(data)
    #print(graphs)
    #data = {(7, 6): [{'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_37', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 268, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_36', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 269, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_31', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 277, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_34', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 269, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}], (7, 5): [{'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_14', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 270, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_11', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 291, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_10', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 292, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 7, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_16', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 309, '# Total False Positives': 0, 'Detection Time': 255.5, '# True Positive Rejections': 57, '# False Negatives': 0, '# False Positive Rejections': 0, '# False Positives': 0}], (6, 5): [{'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_4', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 967, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_9', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 1002, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_8', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 988, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_5', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 1030, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_7', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 981, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_1', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 5, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 980, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}], (6, 6): [{'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_28', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 979, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_25', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 987, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_27', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 983, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_24', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 997, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_29', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 963, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}, {'config': {'sensorSamplingPeriod': 1000, 'numSuperNodes': 4, 'serverSamplingPeriod': 1000, 'GridCapacityPercentage': 0.9, 'detectionDistance': 6, 'stimFileName': 'circle_0.txt', 'movementSamplingSpeed': 20, 'doOptimize': True, 'csvMove': True, 'memprofile': '', 'GPSSamplingPeriod': 1000, 'logGrid': True, 'SamplingLossAccelCM': 0.001, 'OutputFileName': '/home/simulator/bigData/Log_26', 'StandardDeviationThreshold': 1.7, 'movementPath': 'marathon_street_2k.txt', 'SamplingLoss4GCM': 0.005, 'movementSamplingPeriod': 20, 'detectionThreshold': 6, 'GPSSamplingLoss': 0.005, 'thresholdBatteryToUse': 10, 'inputFileName': 'Scenario_3.txt', 'serverSamplingLoss': 0.01, 'energyModel': 'variable', 'SuperNodeSpeed': 3, 'nodeStoredSamples': 10, 'noEnergy': True, 'regionRouting': True, 'negativeSittingStopThreshold': -10, 'sittingStopThreshold': 5, 'csvSensor': True, 'SquareRowCM': 60, 'errorMultiplier': 1, 'cpuprofile': '', 'sensorSamplingLoss': 0.01, 'naturalLoss': 0.005, 'superNodes': True, 'SquareColCM': 320, 'GridStoredSamples': 10, 'logEnergy': True, 'iterations': 999, 'SamplingLossBTCM': 0.0001, 'maxBufferCapacity': 25, 'logNodes': True, 'thresholdBatteryToHave': 30, 'logPosition': True, 'SamplingLossWifiCM': 0.001, 'validationThreshold': 1, 'imageFileName': 'marathon_street_map.png', 'sensorPath': 'smoothed_marathon.csv', 'Recalibration Threshold': 3, 'outRoutingStatsName': 'routingStats.txt'}, '# Total False Negatives': 950, '# Total False Positives': 0, 'Detection Time': 271.5, '# True Positive Rejections': 87, '# False Negatives': 1, '# False Positive Rejections': 0, '# False Positives': 0}]}
    #graphs = [[{'~': 'detectionThreshold', 'detectionDistance': 6}, {(6, 5), (6, 6)}], [{'~': 'detectionThreshold', 'detectionDistance': 7}, {(7, 6), (7, 5)}], [{'~': 'detectionDistance', 'detectionThreshold': 5}, {(7, 5), (6, 5)}], [{'~': 'detectionDistance', 'detectionThreshold': 6}, {(7, 6), (6, 6)}]]

    y_axes = [x for x in data[list(data.keys())[0]][0].keys() if x is not 'config']
    print(y_axes)

    print(len(graphs))

    for graphSet in graphs:
        d = {}
        x_axis = graphSet[0]['~']
        extension = ''
        for k in graphSet[0].keys():
            if k is not '~':
                extension+= str(k) + str(graphSet[0][k])
        for key in graphSet[1]:
            v = data[key][0]['config'][x_axis]

            if v in d:
                d[v].append(data[key])
            else:
                d[v] = data[key]

        print(d.keys())
        print('---')

        for g in y_axes:
            datum = {}#[0 for i in d.keys()]
            #count = #[0 for i in d.keys()]

            for k in d.keys():
                #print(k)
                dat = [x[g] for x in d[k]]
                datum[k] = sum(dat)/len(dat)

            x_val = sorted(list(datum.keys()))
            y_val = [datum[x] for x in x_val]

            plt.plot(x_val, y_val)
            plt.title(g)
            plt.xlabel(x_axis)
            plt.ylabel(g)
            plt.savefig('%s(%s) %s %s.png' % (figurePath, extension, x_axis, g))
            plt.clf()
            #plt.show()

    




















