from zipfile import *
import os

def isInt(s):
    try:
        int(s)
        return True
    except ValueError:
        return False

def isFloat(s):
    try:
        float(s)
        return True
    except ValueError:
        return False

def isBool(s):
    return (s=="true\n" or s=="false\n")

def getParameters(basename):
    zf = ZipFile(basename)
    temp = os.path.split(basename)[1]
    n = temp.split(".zip")[0]

    if 'p2' in n:
        n = n[:-2]

    f = zf.open("%s%s" % (n, "-parameters.txt"))
    params = dict()

    val = ''
    for line in f:
        line = line.decode("utf-8")
        l=line.rstrip().split('=')
        key = l[0]
        if(isBool(l[1])):
            if l[1] == 'true':
                val = True
            else:
                val = False
        elif(isInt(l[1])):
            val = int(l[1])
        elif(isFloat(l[1])):
            val = float(l[1])
        else:
            val = l[1].split("\n")[0]
        params[key] = val
    return params


