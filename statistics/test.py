import pickle

f = open('cluster_data.pickle', 'rb')

data = pickle.load(f)

print(len(data))
#print(data)