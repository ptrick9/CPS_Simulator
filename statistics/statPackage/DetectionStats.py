import re
import numpy as np
import math
from zipfile import *
import os


def dist(x1, y1, x2, y2):
    return math.sqrt((x2-x1)**2 + (y2-y1)**2)

Nodes = {}


class Detection:
    def __init__(self, ident=-1, time=-1, status='', conf='', need=-1, typ='', dist=-1, cleanADC=-1, errorADC=-1, senseError=-1, senseClean=-1, rawConc=-1, cause=''):
        self.ident = ident
        self.time = time
        self.status = status
        self.conf = conf
        self.need = need
        self.typ = typ
        self.dist = dist
        self.cleanADC = cleanADC
        self.errorADC = errorADC
        self.senseError = senseError
        self.senseClean = senseClean
        self.rawConc = rawConc
        self.cause = cause

        if self.status == 'Rejection' and self.typ == 'TP':
            if self.ident in Nodes:
                if self.time - Nodes[self.ident] > 60:
                    Nodes[self.ident] = self.time
                    self.typ = 'FN'
                    print("changed")
                else:
                    pass
            else:
                Nodes[self.ident] = self.time
                self.typ = 'FN'
                print("changed")

    def __lt__(self, other):
        return self.time < other.time

    def __repr__(self):
        return "%d %s %s %f %f" %(self.ident, self.typ, self.status, self.time, self.dist)

    def TPConf(self):
        return self.typ == 'TP' and self.status == 'Confirmation'

    def TPRej(self):
        return self.typ == 'TP' and self.status == 'Rejection'

    def FPConf(self):
        return self.typ == 'FP' and self.status == 'Confirmation'

    def FPRej(self):
        return self.typ == 'FP' and self.status == 'Rejection'

    def FN(self):
        return self.typ == 'FN'

    def FP(self):
        return self.typ == 'FP'

    def Drift(self):
        return self.cause == ' Drift '

    def Wind(self):
        return self.cause == ' Wind '
    #def FNConf(self):
    #    return self.typ == 'FN' and self.status == 'Confirmation'

    #def FNRej(self):
    #    return self.typ == 'FN' and self.status == 'Rejection'


def CountFPReject(d):
    return sum([1 if x.FPRej() else 0 for x in d])

def CountFPConfirmation(d):
    return sum([1 if x.FPConf() else 0 for x in d])

def CountTPReject(d):
    return sum([1 if x.TPRej() else 0 for x in d])

def CountTPConfirmation(d):
    return sum([1 if x.TPConf() else 0 for x in d])

def CountFN(d):
    return sum([1 if x.FN() else 0 for x in d])

def BuildDetections(basename):

    zf = ZipFile(basename)
    temp = os.path.split(basename)[1]
    n = temp.split(".zip")[0]

    if 'low' in n:
        n = n[:-3]

    f = zf.open("%s%s" % (n, "-detection.txt"))


    #f = open('%s-detection.txt' % basename)

    longBoi = ""

    lines = []
    for line in f:
        line = line.decode("utf-8")
        #longBoi += line
        lines.append(line)

    longBoi = ' '.join(lines)
    detections = []
    initialDetections = {}

    regexConf = r"(?P<status>Rejection|Confirmation) T: (?P<time>\d+) ID: (?P<ident>\d+) (?P<conf>\d+)\/(?P<need>\d+)"
    regexDetail = r"(?P<type>[T,F][P,N])(?P<cause> Wind | | Drift )T: (?P<time>\d+) ID: (?P<ident>\d+) .* D: (?P<distance>\d*\.?\d+) " \
                  r"C: (?P<clean>\d+) E: (?P<error>\d+) SE: (?P<errorSense>\d+.\d+) S: (?P<cleanSense>\d+.\d+) R: (?P<raw>\d+.\d+)"

    detailedMatches = re.finditer(regexDetail, longBoi, re.MULTILINE)

    for matchNum, match in enumerate(detailedMatches, start=1):
        if match.group('type') == 'FN':
            ident = int(match.group('ident'))
            time = float(match.group('time'))/1000
            typ = match.group('type')
            dist = float(match.group('distance'))
            cleanADC = int(match.group('clean'))
            errorADC = int(match.group('error'))
            senseError = float(match.group('errorSense'))
            senseClean = float(match.group('cleanSense'))
            rawConc = float(match.group('raw'))
            cause = match.group('cause')
            detections.append(Detection(ident, time, '', -1, -1, typ, dist, cleanADC, errorADC, senseError, senseClean, rawConc, cause))
        initialDetections[(match.group('ident'), match.group('time'))] = match


    confMatches = re.finditer(regexConf, longBoi, re.MULTILINE)

    for matchNum, match in enumerate(confMatches, start=1):
        try:
            detail = initialDetections[(match.group('ident'), match.group('time'))]
            ident = int(detail.group('ident'))
            time = float(detail.group('time'))/1000
            status = match.group('status')
            conf = int(match.group('conf'))
            need = int(match.group('need'))
            typ = detail.group('type')
            dist = float(detail.group('distance'))
            cleanADC = int(detail.group('clean'))
            errorADC = int(detail.group('error'))
            senseError = float(detail.group('errorSense'))
            senseClean = float(detail.group('cleanSense'))
            rawConc = float(detail.group('raw'))
            detections.append(Detection(ident, time, status, conf, need, typ, dist, cleanADC, errorADC, senseError, senseClean, rawConc))
        except:
            print(basename, (match.group('ident'), match.group('time')))

    detections = sorted(detections)
    return detections


def buildApproachDistances(basename, maxDistance, granularity):

    zf = ZipFile(basename)
    temp = os.path.split(basename)[1]
    n = temp.split(".zip")[0]
    f = zf.open("%s%s" % (n, "-simulatorOutput.txt"))

    #f = open('%s-simulatorOutput.txt' % basename)
    lines = []
    for line in f:
        line = line.decode("utf-8")
        lines.append(line.rstrip())

    bx = int(lines[4].split(" ")[2])
    by = int(lines[5].split(" ")[2])
    print(bx, by)

    regexTime = r"t=[ ]*(?P<time>\d+)[ ]*amount=[ ]*(?P<amount>\d+)"
    regexInfo = r"ID:[ ]*(?P<id>\d+)[ ]*x:[ ]*(?P<x>\d+)[ ]*y:[ ]*(?P<y>\d+)"

    approachTime = {}


    maxInds = int(maxDistance/granularity)

    time = 0
    for line in lines:
        if 't=' in line:
            match = [m.groupdict() for m in re.finditer(regexTime, line)]
            time = int(match[0]['time'])
        if 'ID' in line:
            match = [m.groupdict() for m in re.finditer(regexInfo, line)]
            ident = int(match[0]['id'])
            x = int(match[0]['x'])
            y = int(match[0]['y'])
            d = dist(x, y, bx, by)/2.0
            ind = int(d/granularity)
            if ident in approachTime:
                if ind < maxInds:
                    if time < approachTime[ident][ind]:
                        approachTime[ident][ind] = time
                        #print("updated %d %d %d %d\n" % (ident, ind*granularity, d, time))
            else:
                dists = [1000000 for i in range(maxInds)]
                if ind < maxInds:
                    dists[ind] = time
                approachTime[ident] = dists


    v = list(approachTime.values())
    approach = [sorted(v, key= lambda times: times[x]) for x in range(maxInds)]

    a = [[x[i] for x in approach[i]] for i in range(maxInds)]








"""
app = np.asarray(a)

times = [i for i in range(10,999,10)]
dd = [i for i in range(0, 100, 5)]
counts = []
for t in times:
    time_counts = [sum(d < t) for d in app]
    counts.append(time_counts)

times = [[i for x in range(20)] for i in range(10,999,10)]
dd = [[i for i in range(0, 100, 5)] for x in range(99)]

times = np.asarray(times)
dd = np.asarray(dd)
counts = np.asarray(counts)

print(times.shape)
print(dd.shape)
print(counts.shape)
"""