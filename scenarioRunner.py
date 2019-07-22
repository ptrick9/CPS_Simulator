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


    scenarios = ["-inputFileName=%s -imageFileName=%s -logPostion=true -logGrid=false -logEnergy=false -logNodes=false" \
                 "-noEnergy=true -csvMove=true -iterations=1000 -superNodes=false" % (s[0], s[1]) for s in [['Scenario_3.txt', 'marathon_street_map.png']]]

    row = ["-squareRow=%d" % d for d in [60, 120]]
    col = ["-squareCXol=%d" % d for d in [320, 640]]
    movementPath = ["-movementPath=%s" % s for s in ["marathon_street_2k.txt"]]
    detectionThreshold = ["-detectionThreshold=%d" % d for d in[5, 6, 7]]
    detectionDistance = ["-detectionDistance=%d" % d for d in [6, 7]]


    runs = (list(itertools.product(*[scenarios, row, col, movementPath, detectionThreshold, detectionDistance])))
    
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
