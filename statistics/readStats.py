from statPackage.DetectionStats import *

'''
import pickle
data = {}
with open('driftExplorePar.pickle', 'rb') as handle:
    #print(pickle.load(handle))
    data = pickle.load(handle)



i = 0
for k in data.keys():
    i += len(data[k])
print(i)

'''

basePath = "C:/Users/patrick/Downloads/driftExploreCommBombHull/"
file = 'Log_1221.zip'



def getDetections(basePath, run):
    return BuildDetections("%s%s" % (basePath, run))


det = getDetections(basePath, file)

run = {}

firstDet = 20000

run['# False Negatives'] = sum([1 if x.FN() and x.time < firstDet else 0 for x in det])
run['# False Negatives Drift'] = sum([1 if x.FN() and x.Drift() and x.time < firstDet else 0 for x in det])


for x in det:
    print(x.cause)


print(run)



