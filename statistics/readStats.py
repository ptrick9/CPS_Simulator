import pickle
data = {}
with open('driftExplorePar.pickle', 'rb') as handle:
    #print(pickle.load(handle))
    data = pickle.load(handle)

print(data)
y_axes = [x for x in data[list(data.keys())[0]][0].keys() if x is not 'config']
print(y_axes)