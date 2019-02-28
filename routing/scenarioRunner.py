import multiprocessing
import itertools
import os
import random
import subprocess as sp
import time



def runner(queue):
    while True:
        job = queue.get()
        command = "./src "+' '.join(job)
        #print(command)
        s = time.time()
        FNULL = open(os.devnull, 'w')
        job.insert(0, "./routing")
        child = sp.Popen(job, stdout=FNULL, stderr=sp.PIPE)
        child.wait()
        e = time.time()
        FNULL.close()
        print("%f for %s" % (e-s, command))
        queue.task_done()


if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()

    num_nodes = ["-numSuperNodes=%d" % i for i in [1, 2, 4]]
    scenario = ["circle_justWalls_x4", "maze_justWalls_x4", "school_justWalls_x4"]

    runMeat = []
    for s in scenario:
      for i in range(1, 40):
        runMeat.append("-imageFileName=%s.png -stimFileName=routing-test-traces/%s_%d.txt -outRoutingStatsName=stats/%s_%d_stats.txt" % (s, s, i, s, i))
    
    routeType = ["-regionRouting", "-aStarRouting"]

    vals = list(itertools.product(*[num_nodes, runMeat, routeType]))
    print(vals)
    


    
    #print(runs)
    
    x = 0
    for r in runs:
        #paramsInter = [rotateFrequnecy, rotationFactorInter]
        #interRuns = (list(itertools.product(*paramsInter)))
        #interRunsDone = []
        for i in range(10):
            #v = list(itertools.chain(r, ["-deadlineOutputFile=Log_%d" % x]))
            q.put(v)
            x+= 1
    
        
    p = multiprocessing.Pool(24, runner, (q,))

    q.join()
