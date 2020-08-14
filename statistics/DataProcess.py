import matplotlib.pyplot as plt
import numpy as np

    


def BarGraphs():
    """
    Follow the comments in the first few lines
    """
    spacing=np.arange(40,1080,80)
    #Spacing is the iteration range you are looking at
    #spacing=np.arange(40,320,20)
    valueList=[]
    #List of values to plot
    maxval=140
    #Maxval is for old functionlity that can be removed with no effect to the code
    fig, ax = plt.subplots(3,2)
    #The with loop below opens whatever file you are reading from, this is the file you write to using LogReader.py
    with open(r'C:\Users\brook\Desktop\Data Files\3DDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    
    ax[0,0].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    #Using ax[0,0] plots at the sub plot 0,0 and you use clipping to make sure no values are excluded
    for rect in ax[0,0].patches:
        height = rect.get_height()
        ax[0,0].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    #The code above puts the numbers above the bar graphs so you can read the values
    ax[0,0].set_xticks(spacing)
    ax[0,0].set_title("3x3 Adaptive Deltas")
    ax[0,0].set_xlabel("Time Threshold")
    ax[0,0].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    with open(r'C:\Users\brook\Desktop\Data Files\5DDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    ax[1,0].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    for rect in ax[1,0].patches:
        height = rect.get_height()
        ax[1,0].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    ax[1,0].set_xticks(spacing)
    ax[1,0].set_title("5x5 Adaptive Deltas")
    ax[1,0].set_xlabel("Time Threshold")
    ax[1,0].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    with open(r'C:\Users\brook\Desktop\Data Files\8DDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    ax[2,0].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    for rect in ax[2,0].patches:
        height = rect.get_height()
        ax[2,0].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    ax[2,0].set_xticks(spacing)
    ax[2,0].set_title("8x8 Adaptive Deltas")
    ax[2,0].set_xlabel("Time Threshold")
    ax[2,0].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    with open(r'C:\Users\brook\Desktop\Data Files\3MDDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    ax[0,1].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    for rect in ax[0,1].patches:
        height = rect.get_height()
        x=0
        if height> 800:
            x=800
        ax[0,1].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height-x), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    ax[0,1].set_xticks(spacing)
    ax[0,1].set_title("3x3 Adaptive Deltas Max")
    ax[0,1].set_xlabel("Time Threshold")
    ax[0,1].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    with open(r'C:\Users\brook\Desktop\Data Files\5MDDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    ax[1,1].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    for rect in ax[1,1].patches:
        height = rect.get_height()
        x=0
        if height> 800:
            x=400
        ax[1,1].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height-x), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    ax[1,1].set_xticks(spacing)
    ax[1,1].set_title("5x5 Adaptive Deltas Max")
    ax[1,1].set_xlabel("Time Threshold")
    ax[1,1].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    with open(r'C:\Users\brook\Desktop\Data Files\8MDDAS.txt') as f:
        for line in f:
            data=line.split()
            valueList.append(float(data[0]))
            if float(data[0]) > maxval:
                maxval=float(data[0])
    ax[2,1].hist(np.clip(valueList,spacing[0],spacing[-1]), density=False,bins=spacing,)
    for rect in ax[2,1].patches:
        height = rect.get_height()
        x=0
        if height> 800:
            x=200
        ax[2,1].annotate(f'{int(height)}', xy=(rect.get_x()+rect.get_width()/2, height-x), 
                    xytext=(0, 5), textcoords='offset points', ha='center', va='bottom') 
    ax[2,1].set_xticks(spacing)
    ax[2,1].set_title("8x8 Adaptive Deltas Max")
    ax[2,1].set_xlabel("Time Threshold")
    ax[2,1].set_ylabel("# of Delta Occurences")
    valueList.clear()
    f.close()
    plt.show()

def showBatteryStats():
    """
    This function is used to generate one plot with 4 graphs
    to compare different algorithms effectivness
    """
    fig, ax = plt.subplots()
    SamplesBO=[3306588/3500, 3355327/3500, 2162116/3500, 2204989/3500]
    SamplesBT=[2765658/3500, 3090913/3500, 2043365/3500, 2102971/3500]
    NodeNum=3500
    vals=[0,1,2,3]
    names=["Non Adaptive","Speed Adaptation", "Density Adaptation", "Both Adaptations"]
    ax.bar(vals,height=SamplesBT,width=.4)
    plt.xticks(vals,names)
    plt.title("Samples per node with Battery Adaptation")
    for rect in ax.patches:
        height = rect.get_height()
        ax.text(rect.get_x() + rect.get_width()/2., height,
            '%.2f' % height, ha='center', va='bottom')

    plt.show()

def StackedBarGraph():
    header = ['0 Samples','1-2 Samples','3-10 Samples','11+ Samples']
    ind=[0,1,2,3]
    data=[[3055, 59, 626, 8493],[3073, 42, 545, 8579],[3062, 42, 484, 8651],[3058, 41, 452, 8694]]
    zeroes=[]
    var1=[]
    var2=[]
    var3=[]
    for i in range(len(data)):
        zeroes.append(data[i][0])
        var1.append(data[i][1])
        var2.append(data[i][2])
        var3.append(data[i][3])
    p1=plt.bar(ind,zeroes,color='r',width=.25)
    p2=plt.bar(ind,var1,bottom=np.array(zeroes),color='b',width=.25)
    p3=plt.bar(ind,var2,bottom=np.array(zeroes)+np.array(var1),color='g',width=.25)
    p4=plt.bar(ind,var3,bottom=np.array(zeroes)+np.array(var1)+np.array(var2),color='c',width=.25)
    plt.legend((p1[0], p2[0], p3[0], p4[0]), (header[0], header[1], header[2], header[3]), fontsize=10, ncol=4, framealpha=0, fancybox=True,loc="upper right")
    names=["Non Adaptive", "Speed Adaptation", "Density Adpatation", "Both Adaptations"]
    plt.xticks(ind,names)
    plt.title("Samples per Grid Square with Battery Adaptation")
    plt.show()
