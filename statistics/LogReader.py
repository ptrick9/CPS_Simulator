maxDeltas=[]
Deltas=[]
coverage=[]
hasWritten=False
#This with loop reads the data log and can be configured however to withdraw the data from it
with open(r"C:\Users\brook\Desktop\OutputFolder\test-OutputLog.txt","r") as f:
        data=f.readlines()
        for line in data:
            line=line.split()
            if line[2]=="max":
                maxDeltas.append(str(round(float(line[5])/1000,1)))
                coverage.append(float(line[6]))
            else:
                Deltas.append(str(round(float(line[5])/1000,1)))
f.close()

def write():
    """
    Writes to specified files the data you want from the output log
    """
    with open(r"C:\Users\brook\Desktop\Data Files\3DDAS.txt","w") as f:
        for items in range(len(Deltas)):
            f.write(Deltas[items]+"\n")
    f.close()

    with open(r"C:\Users\brook\Desktop\Data Files\3MDDAS.txt","w") as f:
        for items in range(len(maxDeltas)):
            f.write(maxDeltas[items]+"\n")
    f.close()

    global hasWritten
    hasWritten=True
def calcCoverage():
    print("Coverage :",len(coverage))
    numAccGrids=12274
    var1=numAccGrids-len(coverage)
    var2=0
    var3=0
    var4=0
    for i in range(len(coverage)):
        if coverage[i] >1 and coverage[i] < 3:
            var2+=1
        elif coverage[i] >=3 and coverage[i] <11:
            var3+=1
        elif coverage[i]>=11:
            var4+=1
    print(var1,var2,var3,var4)
#write()
#calcCoverage
#This is a check to make sure you are writting to the file
if hasWritten==True:
    print("Successfully wrote to files this time")
else:
    print("Have not written to files this time")
print("Deltas:",len(Deltas))
print("MaxDeltas :",len(maxDeltas))

