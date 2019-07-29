from statPackage.DetectionStats import *
from statPackage.ParamProcessing import *
import os
import itertools

#basePath = "C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/data/bigData/"
basePath = "C:/Users/patrick/Downloads/bigData/"

X_VAL = 'validationThreshold'


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
            if 'file' not in k and 'File' not in k:
                if k in params:
                    params[k].add(p[k])
                else:
                    params[k] = set()
                    params[k].add(p[k])

    unique = {}
    for k in params.keys():
        if len(params[k]) > 1:
            unique[k] = list(sorted(params[k]))
    print(unique)
    return unique




uniqueRuns = set()
for file in os.listdir(basePath):
    uniqueRuns.add(file.split('-')[0])



shiftingParameters = determineDifferences(basePath, uniqueRuns)

vals = []
order = []
for k in shiftingParameters.keys():
    order.append(k)
    vals.append(shiftingParameters[k])

print(order)
print(list(itertools.product(*vals)))

#processed =

data = {}

for i,file in enumerate(uniqueRuns):
    run = {}
    p = getParameters("%s%s" % (basePath, file))
    det = getDetections(basePath, file)
    firstDet = next(x for x in det if x.TPConf()).time
    run['Detection Time'] = firstDet
    run['# True Positive Rejections'] = sum([1 if x.TPRej() and x.time < firstDet else 0 for x in det])
    run['# False Positives'] = sum([1 if x.FP() and x.time < firstDet else 0 for x in det])
    run['# False Positive Rejections'] = sum([1 if x.FPRej() and x.time < firstDet else 0 for x in det])
    run['# False Negatives'] = sum([1 if x.FN() and x.time < firstDet else 0 for x in det])
    run['# Total False Negatives'] = sum([1 if x.FN() else 0 for x in det])
    run['# Total False Positives'] = sum([1 if x.FP() else 0 for x in det])
#k = tuple(p[x] for x in order)
    k = p[X_VAL]
    if k in data:
        data[k].append(run)
    else:
        data[k] = [run]

    #processed
    #print(run)
    print("%d\%d" % (i, len(uniqueRuns)))
print(data)





#buildDetectionList(basePath, uniqueRuns)



"""
for i,val in enumerate(range(0, maxDistance, granularity)):
    print("Earliest <%dm: %d" % (val, a[i][0]))



print(CountFPReject(detections))
print(CountFPConfirmation(detections))
print(CountTPReject(detections))
print(CountTPConfirmation(detections))
print(CountFN(detections))


print(next(x for x in detections if x.TPRej()))
print(next(x for x in detections if x.TPConf()))
"""

