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
    filename = "%s-parameters.txt" % basename
    params = dict()
    with open(filename) as f:
        lines = f.readlines()

    for i in range(0,len(lines)-1):
        l=lines[i].split('=')
        key = l[0]
        if(isBool(l[1])):
            val = bool(l[1])
        else:
            if(isInt(l[1])):
                val = int(l[1])
            else:
                if(isFloat(l[1])):
                    val = float(l[1])
                else:
                    val = l[1].split("\n")[0]
        params[key] = val
    return params

