from statPackage.DetectionStats import *
from statPackage.ParamProcessing import *
import multiprocessing


import os
import itertools
import matplotlib.pyplot as plt
import pickle
from zipfile import *

#basePath = "C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/data/bigData/"
#basePath = "C:/Users/patrick/Downloads/bigData/"
#basePath = "C:/Users/patrick/Downloads/fineGrainedBomb/fineGrainedBomb/"
#basePath = "C:/Users/patrick/Downloads/driftExploreHullBombMove/"
basePath = "C:/Users/patrick/Downloads/driftExplorerNoBomb/"
#basePath = "C:/Users/patrick/Downloads/driftTest/"
figurePath = "C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2018/git_simulator/CPS_Simulator/driftExploreCommBomb/"

X_VAL = ['detectionThreshold', 'detectionDistance']
X_VAL = ['validationThreshold', 'errorMultiplier', 'serverRecal']

IGNORE = ['movementPath', 'bombX', 'bombY']
ZIP = True

data = {}


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
        print(r)
        print(p['serverRecal'])
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

def runner(queue):
    while True:
        job = queue.get()
        print("%d\\%d" % (job[1], job[2]))
        #print(job)
        command = "./simulator/Simulation "+' '.join(job[0])
        print(command)
        #os.system(command + " 1>/dev/null")
        queue.task_done()

# Generation of Data
def generateData(rq, wq):

    while True:
        job = rq.get()
        file = job[0]
        order = job[1]
    #data = {}
    #allData = {}
        #for i,file in enumerate(uniqueRuns):
        run = {}
        print(job)
        p = getParameters("%s%s" % (basePath, file))
        det = getDetections(basePath, file)
        firstDet = 20000
        try:
            firstDet = next(x for x in det if x.TPConf()).time
        except:
            pass
        run['Detection Time'] = firstDet
        run['# True Positive Rejections'] = sum([1 if x.TPRej() and x.time < firstDet else 0 for x in det])
        run['# False Positives'] = sum([1 if x.FP() and x.time < firstDet else 0 for x in det])
        run['# False Positive Confirmations'] = sum([1 if x.FPConf() and x.time < firstDet else 0 for x in det])
        run['# False Positive Rejections'] = sum([1 if x.FPRej() and x.time < firstDet else 0 for x in det])
        run['# False Positive Wind'] = sum([1 if x.FP() and x.time < firstDet and x.Wind() else 0 for x in det])
        run['# False Positive Drift'] = sum([1 if x.FP() and x.time < firstDet and x.Drift() else 0 for x in det])
        run['# False Negatives'] = sum([1 if x.FN() and x.time < firstDet else 0 for x in det])
        run['# False Negatives Drift'] = sum([1 if x.FN() and x.Drift() and x.time < firstDet else 0 for x in det])
        run['# Total False Negatives'] = sum([1 if x.FN() else 0 for x in det])
        run['# Total False Positives'] = sum([1 if x.FP() else 0 for x in det])

        run['config'] = {}
        for k in p.keys():
            run['config'][k] = p[k]

        k = tuple(p[x] for x in order)

        wq.put([run, det, k])

        '''
        #k = p[X_VAL]
        if k in data:
            data[k].append(run)
            allData[k].append(det)
        else:
            data[k] = [run]
            allData[k] = [det]
        

        print("\r%d\%d" % (i, len(uniqueRuns)), end='')
        print('')
        print(data)
        with open('driftExplore.pickle', 'wb') as handle:
            pickle.dump({'det': allData, 'processed': data}, handle)# protocol=pickle.HIGHEST_PROTOCOL)
            #pickle.dump({data, handle)
        return data
        '''
        rq.task_done()

def dataStorage(wq, order, variation, total):

    i = 0
    failed = 0
    while True:
        data = {}
        job = wq.get()
        try:
            with open('driftExploreNoBomb23.pickle', 'rb') as handle:
                data = pickle.load(handle)# protocol=pickle.HIGHEST_PROTOCOL)
                data = data['data']
                handle.close()
        except:
            pass
        run = job[0]
        det = job[1]
        k = job[2]

        if k in data:
            data[k].append(run)
            #allData[k].append(det)
        else:
            data[k] = [run]
            #allData[k] = [det]

        dat = {'data': data, 'order': order, 'var': variation}
        fail = True
        try:
            f = open('driftExploreNoBomb23.pickle', 'wb')
            pickle.dump(dat, f)
            f.close()
            fail = False
        except:
            failed += 1
            pass
        i += 1
        print("Total %d/%d Failed %d/%d" % (i, total, failed, i))
        #with open('driftExploreBomb.pickle', 'wb') as handle:
        #    pickle.dump(data, handle)# protocol=pickle.HIGHEST_PROTOCOL)
        #    handle.close()
        wq.task_done()

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
    uniqueRuns = list(uniqueRuns)#[:100]

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

    m = multiprocessing
    rq = m.JoinableQueue()
    wq = m.JoinableQueue()

    for run in uniqueRuns:
        rq.put([run, order])

    p = multiprocessing.Pool(3, generateData, (rq,wq,))
    p = multiprocessing.Pool(1, dataStorage, (wq,order, shiftingParameters, len(uniqueRuns),))

    rq.join()



    wq.join()

    #with open('driftExplorePar.pickle', 'wb') as handle:
    #    pickle.dump(data, handle)# protocol=pickle.HIGHEST_PROTOCOL)
    #    handle.close()
        #pickle.dump({data, handle)

    #data = generateData(uniqueRuns)

    with open('driftExplorePar.pickle', 'rb') as handle:
        data = pickle.load(handle)

    print(data)
    graphs = generateGraphs(order, xx)

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

    




















