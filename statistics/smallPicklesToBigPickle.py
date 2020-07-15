import os
import pickle

def smallPicklesToBigPickle(numberOfRuns, pickleName):

    i = 1
    failed = 0
    pickleOpen = False
    data = {}
    allData = {}
    waiting = 0
    while i <= numberOfRuns:
        try:
            with open('%s%d.pickle' % (pickleName, i), 'rb') as handle:
                runData = pickle.load(handle)# protocol=pickle.HIGHEST_PROTOCOL)
                run = runData['run']
                k = runData['key']
                if k in data:
                    data[k].append(run)
                else:
                    data[k] = [run]
                handle.close()
            os.remove('%s%d.pickle' % (pickleName, i))
        except:
            failed += 1
        print("wq:", i, failed)
        i += 1

    allData = {'data': data, 'order': runData['order'], 'var': runData['var']}
    try:
        f = open('%s.pickle' % (pickleName), 'wb')
        pickle.dump(allData, f)
        f.close()
    except:
        pass

    return data


smallPicklesToBigPickle(1552, 'LRTests')