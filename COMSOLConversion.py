import csv
import numpy as np
from scipy.ndimage import gaussian_filter


file_name = 'C:/Users/patrick/Dropbox/CPS_EXPLOSIVES/SimulationData/Marathon/Marathon_2D_10SecR_1000s.csv'

header = ''
f = open(file_name, 'r')
for line_number, line in enumerate(f):
    if line_number == 8:
        header = line.rstrip()
    elif line_number > 8:
        break

new_header = ['X', 'Y']
new_header.extend(list(map(lambda x: x.split('@ t=')[1], header.split(',')[3:])))
print(new_header)

jj = np.genfromtxt(file_name, delimiter=',', skip_header=9)

x = abs(jj[:,0])
y = abs(jj[:,1])


width = int(max(x)*2+1)
height = int(max(y)*2+1)


data = np.zeros((width, height, len(jj[0])-3))

for i,coord in enumerate(zip(x, y)):
    data[abs(int(coord[0]*2))][height-abs(int(coord[1]*2))-1] = jj[i][3:]


t = data[:,:,4]
s = t.shape
m = 0
mx = 0
my = 0
for i in range(s[0]):
    for j in range(s[1]):
        if t[i][j] > m:
            m = t[i][j]
            mx = i
            my = j
print((mx)/2, my/2)

averaged = np.zeros((width, height, len(jj[0])-3))
for t in range(len(jj[0])-3):
    time_step = data[:,:,t]
    averaged[:,:,t] = gaussian_filter(time_step, sigma=3)


printWidth = width - width%10 #round down to nearest 10
printHeight = height - height%10 #round down to nearest 10

with open('smoothed_marathon.csv', 'w') as writeFile:
    writer = csv.writer(writeFile, lineterminator='\n')
    writer.writerow(new_header)
    for i in range(printWidth):
        for j in range(printHeight):
            formatted = list(map(lambda t: "%e" % t if t != 0 else "0.0",(list(averaged[i,j,:]))))
            temp = [i, j]
            temp.extend(formatted)
            writer.writerow(temp)