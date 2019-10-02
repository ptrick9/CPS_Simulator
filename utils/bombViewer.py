'''
import math
import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D

def dist(x1, y1, x2, y2):
    return math.sqrt((x2-x1)**2 + (y2-y1)**2)


import csv
import numpy as np
from scipy.ndimage import gaussian_filter, spline_filter, median_filter

import imageio

file_name = 'C:/Users/patrick/Dropbox/CPS_EXPLOSIVES/SimulationData/Marathon/marathon_v2/Marathon_2D_10SecR_5000s_2_NoWind.csv'
#file_name = 'C:/Users/patrick/Dropbox/CPS_EXPLOSIVES/SimulationData/Marathon/marathon_v2/Marathon_2D_10SecR_5000s.csv'
#file_name = 'C:/Users/patrick/Dropbox/CPS_EXPLOSIVES/SimulationData/Simulation_UDel_Geoemtry/Circle_Just_walls_20Sec_Resolution/Circle_2D.csv'

header = ''
f = open(file_name, 'r')
for line_number, line in enumerate(f):
    if line_number == 8:
        header = line.rstrip()
    elif line_number > 8:
        break


new_header = ['Scale', 'X', 'Y']
new_header.extend(list(map(lambda x: x.split('@ t=')[1], header.split(',')[3:])))
print(new_header)

jj = np.genfromtxt(file_name, delimiter=',', skip_header=9)

x = abs(jj[:,0])
y = abs(jj[:,1])

from mpl_toolkits.mplot3d import Axes3D
from matplotlib import cm


pts = []
for i,coord in enumerate(zip(x, y)):
    if coord[0] > 4.8 and coord[0] < 5.2 and coord[1] > 4.8 and coord[1] < 5.2:
        pts.append([coord[0], coord[1], jj[i][450]])
#print(pts)
pts.sort(key=lambda c: c[0])
#print(pts)
xx =  [p[0] for p in pts]
yy = [p[1] for p in pts]
zz = [p[2] for p in pts]



fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

ax.plot_trisurf(xx, yy, zz, cmap='viridis', edgecolor='none')

plt.show()
'''

import numpy as np

f = open('C:/Users/patrick/Dropbox/Patrick/udel/SUMMER2019/GitSimulator/fine_bomb9x9.csv', 'r')
jj = np.genfromtxt(f, delimiter=',', skip_header=1)

print(jj[0])

x = jj[:,1]
y = jj[:,2]


import matplotlib.pyplot as plt
from mpl_toolkits.mplot3d import Axes3D
from matplotlib import cm


pts = []
for i,coord in enumerate(zip(x, y)):
    if coord[0] > 20 and coord[0] < 30 and coord[1] > 20 and coord[1] < 30:
        if coord[0] == 22 and coord[1] == 22:
            #pts.append([22, 22, 20000])
            pts.append([coord[0], coord[1], jj[i][135]])
        else:
            pts.append([coord[0], coord[1], jj[i][135]])



#print(pts)
pts.sort(key=lambda c: c[0])
#print(pts)
xx =  [p[0] for p in pts]
yy = [p[1] for p in pts]
zz = [p[2] for p in pts]



fig = plt.figure()
ax = fig.add_subplot(111, projection='3d')

ax.plot_trisurf(xx, yy, zz, cmap='viridis', edgecolor='none')
#ax.scatter(xx, yy, zz, cmap='viridis', edgecolor='none')

ax.set_xlabel('x')
ax.set_ylabel('y')
plt.show()
