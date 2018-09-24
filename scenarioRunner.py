import multiprocessing
import itertools
import os



def runner(queue):
    while True:
        job = queue.get()
        command = "./simulator "+' '.join(job)
        print(command)
        os.system(command)
        queue.task_done()

if __name__ == '__main__':

    m = multiprocessing
    q = m.JoinableQueue()

    #detect_threshold = ["-detectionThreshold=%f" % float(20.0+0.5*i) for i in range(0, 20, 1)]
    scenarios = ["-inputFileName=%s" % s for s in ['Scenario_Random_100.txt', 'Scenario_Random_200.txt', 
     'Scenario_Random_500.txt', 'Scenario_Random_1000.txt', 
     'Scenario_Random_2000.txt', 'Scenario_Random_5000.txt', 'Scenario_Random_10000.txt']]

    error_mult = ["-errorMultiplier=%f" % s for s in [0.0, 0.1, 1.0]]

    runs = (list(itertools.product(*[error_mult, scenarios])))
    
    x = 200
    for r in runs:
        #paramsInter = [rotateFrequnecy, rotationFactorInter]
        #interRuns = (list(itertools.product(*paramsInter)))
        #interRunsDone = []
        for i in range(5):
            v = [r[0], r[1], "-outputFileName=Log_%d" % x, "-squareRow=40 -squareCol=40"]
            q.put(v)
            x+= 1

       
    p = multiprocessing.Pool(15, runner, (q,))

    q.join()
