import multiprocessing
import itertools
import os
import random
import subprocess as sp
import time



def runner(queue):
    while True:
        job = queue.get()
        command = "./routing "+job
        #print(command)
        s = time.time()
        #FNULL = open(os.devnull, 'w')
        #print(job)
        #job = list(job)
        #job.insert(0, "./routing")
        #print(job)
        os.system(command)
        #child = sp.call(job, shell=True)
        #child.wait()
        #os.run(job)
        e = time.time()
        #FNULL.close()
        print("%f for %s" % (e-s, command))
        queue.task_done()


if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()

    num_nodes = ["-numSuperNodes=%d" % i for i in [1, 2, 4]]
    scenario = ["circle_justWalls_x4", "maze_justWalls_x4", "school_justWalls_x4"]

    runMeat = []
    for s in scenario:
      for i in range(1, 200):
        for n in [1, 2, 4]:
          for r in [1, 0]:
            if r == 1:  
              runMeat.append("-numSuperNodes=%d -regionRouting -imageFileName=%s.png -stimFileName=routing-test-traces2/%s_%d.txt -outRoutingStatsName=./stats2/%s_%d_%d_%s_stats.txt" % (n, s, s, i, s, i, n, 'region'))
            else:
              runMeat.append("-numSuperNodes=%d -aStarRouting -imageFileName=%s.png -stimFileName=routing-test-traces2/%s_%d.txt -outRoutingStatsName=./stats2/%s_%d_%d_%s_stats.txt" % (n, s, s, i, s, i, n, 'astar'))
    
    routeType = ["-regionRouting", "-aStarRouting"]

    #vals = list(itertools.product(*[num_nodes, runMeat, routeType]))
    #print(len(vals))
    #print(vals)



    
    #print(runs)
    
    x = 0
    for v in runMeat:
        #paramsInter = [rotateFrequnecy, rotationFactorInter]
        #interRuns = (list(itertools.product(*paramsInter)))
        #interRunsDone = []
        #for i in range(10):
            #v = list(itertools.chain(r, ["-deadlineOutputFile=Log_%d" % x]))
        #v = list(itertools.chain(v))
        q.put(v)
        x+= 1
    
        
    p = multiprocessing.Pool(28, runner, (q,))

    q.join()
