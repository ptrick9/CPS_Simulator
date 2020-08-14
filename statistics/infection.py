import re
import plotly.graph_objects as go

infection_regex = re.compile(r'Risk: (\d+) Infected: (\d+)')


def load_infections_data(path):
    risks = []
    infections = []
    with open(path) as fp:
        line = fp.readline()
        while line is not None:
            search = re.search(infection_regex, line)
            if search is None:
                print("regex not found - ", line, len(risks))
                break
            else:
                risks.append(int(search.group(1)))
                infections.append(int(search.group(2)))
            line = fp.readline()
    return risks, infections


groups = ['0', '10', '30', '60', '90']

figure = go.Figure()
figure.update_layout(title='Total Infected Nodes vs. Time', xaxis_title='time', yaxis_title='number of nodes')
for i, group in enumerate(groups):
    risks, infects = load_infections_data('./' + group + '/test-infection-stats.txt')
    figure.add_trace(go.Scatter(x=list(range(len(infects))), y=infects, mode='lines',
                                name=group + '% of host and non-host nodes wearing masks'))

figure.show()
