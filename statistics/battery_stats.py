import re
import plotly.graph_objects as go

title_regex = re.compile(r'Capacity: (\d+) SamplesLoss: (\d+) WifiLoss: (\d+) BluetoothLoss: (\d+)')
header_regex = re.compile(r'Amount: (\d+) Samples: (\d+) Wifi: (\d+) Bluetooth: (\d+)')
battery_regex = re.compile(r'battery: (-?\d+)')


# Reads the header of the battery log file
# The header contains information about the loss rates of the simulator
# Stores and returns the loss battery loss information in a dictionary
def read_battery_header(file):
    data = {'samples': [], 'wifi': [], 'bluetooth': [], 'battery': {}, 'loss': {}}

    title_string = file.readline()
    search = re.search(title_regex, title_string)
    if search is None:
        print("title not found - ", title_string)

    capacity = int(search.group(1))
    data['loss'] = {
        'samples': int(search.group(2))/capacity,
        'wifi': int(search.group(3))/capacity,
        'bluetooth': int(search.group(4))/capacity,
    }
    return data


# Reads the individual entries of the battery log file
# The entries track the total number of samples/wifi/bluetooth communications over time
# The battery level of each node is also included
# Stores the entries as a list in each corresponding dictionary
# Returns the total number of entries
def read_battery_entries(file, data):
    time = 0
    header_string = file.readline()
    while header_string:
        search = re.search(header_regex, header_string)
        if search is None:
            print('header not found - ', header_string)
            break

        data['samples'].append(int(search.group(2)))
        data['wifi'].append(int(search.group(3)))
        data['bluetooth'].append(int(search.group(4)))
        data['battery'][time] = []
        amount = int(search.group(1))
        for i in range(amount):
            battery_string = file.readline()
            search = re.search(battery_regex, battery_string)
            if search is None:
                print('battery level not found - ', battery_string)
                break
            battery = int(search.group(1))
            data['battery'][time].append(battery)

        header_string = file.readline()
        time += 1

    return time


# Creates an empty graph with the given title, x axis label, and y axis label
def create_graph(title, x, y):
    figure = go.Figure()
    figure.update_layout(title=title, xaxis_title=x, yaxis_title=y)
    return figure


# Creates a graph from battery log data that shows the number of dead nodes in the simulation over time
def create_dead_graph(data, time):
    entries = []
    for index in range(time):
        total = 0
        for i2 in range(len(data['battery'][index])):
            if data['battery'][index][i2] <= 10:
                total += 1
        entries.append(round((total/3500) * 100))
    figure = create_graph('Dead Nodes', 'time', '% of nodes dead')
    figure.add_trace(go.Scatter(x=list(range(len(entries))), y=entries, mode='lines'))
    figure.show()


# Creates a graph from battery log data that shows the change in between
def create_delta_graph(title, x, y, section, data, time):
    entries = []
    index = 10
    while index < time:
        value = data[section][index] - data[section][index - 10]
        entries.append(value)
        index += 10

    figure = create_graph(title, x, y)
    figure.add_trace(go.Scatter(x=[(x + 1) * 10 for x in range(len(entries))], y=entries, mode='lines'))
    figure.show()


# Reads the given battery log file
# Returns a dictionary with the log data and the total time of the simulation run
def read_battery_log(path):
    with open(path) as fp:
        data = read_battery_header(fp)
        time = read_battery_entries(fp, data)
        create_delta_graph('Total Samples Taken', 'time', '# of samples taken', 'samples', data, time)
        create_delta_graph('Total Bluetooth Communications', 'time', '# of bluetooth communications', 'bluetooth', data, time)
        create_delta_graph('Total Wifi Communications', 'time', '# of wifi communications', 'wifi', data, time)


if __name__ == "__main__":
    read_battery_log(r'C:\Users\brook\Desktop\OutputFolder\test-node.txt')