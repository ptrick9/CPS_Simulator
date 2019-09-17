import os

basePath = "C:/Users/patrick/Downloads/driftTest/"


for f in os.listdir(basePath):
    t = f.split('.zip')[0]
    t += 'low'
    newF = '%s%s.zip' % (basePath, t)
    print(f, newF)
    os.rename('%s%s' % (basePath, f), newF)

