import os
import pickle

def smallPicklesToBigPickle(pickleName):

    i = 1
    data = {}
    allData = {}
    while True:
        try:
            with open('%s%d.pickle' % (pickleName, i), 'rb') as handle:
                runData = pickle.load(handle)# protocol=pickle.HIGHEST_PROTOCOL)
                run = runData['run']
                k = runData['key']
                if k in data:
                    data[k] = combineRuns(data[k], run)
                    #data[k].append(run)
                else:
                    data[k] = {'run':run, 'num':1}
                handle.close()
            os.remove('%s%d.pickle' % (pickleName, i))
        except:
            print("end")
            break
        print("wq:", i)
        i += 1

    #data = averageOutRuns(data)

    data = averageCombined(data)

    allData = {'data': data, 'order': runData['order'], 'var': runData['var']}
    try:
        f = open('%s.pickle' % (pickleName), 'wb')
        pickle.dump(allData, f)
        f.close()
    except:
        print("failure")

def averageOutRuns(data):
    for k in data:
        newRun = {}
        i = 0
        for run in data[k]:
            i += 1
            for key in run:
                if key != 'config' and key != 'Distances':
                    if type(run[key]) is list:
                        if key in newRun:
                            listRange = min(len(run[key]), len(newRun[key]))
                            if abs(len(run[key]) - len(newRun[key])) > 2:
                                print(listRange)
                            for j in range(listRange):
                                newRun[key][j] += run[key][j]
                        else:
                            newRun[key] = []
                            for num in run[key]:
                                newRun[key] += [num]
                    else:
                        if key in newRun:
                            newRun[key] += run[key]
                        else:
                            newRun[key] = run[key]
        for key in newRun:
            if type(run[key]) is list:
                for j in range(len(newRun[key])):
                    newRun[key][j] = newRun[key][j]/i
            else:
                newRun[key] += newRun[key]/i
        data[k] = newRun
    return data

def combineRuns(runDict, newRun):
    runDict['num'] += 1
    run = runDict['run']
    for key in run:
        if key != 'config' and key != 'Distances':
            if type(run[key]) is list:
                if key in newRun:
                    listRange = min(len(run[key]), len(newRun[key]))
                    if abs(len(run[key]) - len(newRun[key])) > 2:
                        print(listRange)
                    for j in range(listRange):
                        run[key][j] += newRun[key][j]
            else:
                if key in newRun:
                    run[key] += newRun[key]
    runDict['run'] = run
    return runDict

def averageCombined(data):
    for k in data:
        run = data[k]['run']
        num = data[k]['num']
        for key in run:
            if key != 'config' and key != 'Distances':
                if type(run[key]) is list:
                    for i in range(len(run[key])):
                        run[key][i] = run[key][i]/num
                else:
                    run[key] += run[key]/num
        data[k] = run
    return data

smallPicklesToBigPickle('ReportLocalReclusterNoBattery')