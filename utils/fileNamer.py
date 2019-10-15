import os

basePath = "C:/Users/patrick/Downloads/driftExplorerBombFinal2/"


for f in os.listdir(basePath):
    t = f.split('.zip')[0]
    t += 'p2'
    newF = '%s%s.zip' % (basePath, t)
    print(f, newF)
    os.rename('%s%s' % (basePath, f), newF)

