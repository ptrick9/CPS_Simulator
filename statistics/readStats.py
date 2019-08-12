import pickle
data = {}
with open('driftExplorePar.pickle', 'rb') as handle:
    #print(pickle.load(handle))
    data = pickle.load(handle)



i = 0
for k in data.keys():
    i += len(data[k])
print(i)