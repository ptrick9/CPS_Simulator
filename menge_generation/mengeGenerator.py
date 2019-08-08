import os
import pickle

cwd = os.getcwd()
base_name = 'Basic'
make_graph = False
make_obstacles = True

if make_obstacles:
    import xml.etree.ElementTree as ET
    import xml.dom.minidom as xdm

'''
    Pickle Data Structure
    --------------
    data = {
        'width': int
        'height': int
        'squares': tuple[4]
    }
    
    Tuplet Order for 'squares' in data
    --------------
    x1: Top-Left-X,
    x2: Bottom-Right-X
    y1: Top-Left-Y
    y2: Bottom-Right-Y
'''

# readPickle(string)
# reads a pickle file from the given file
# Returns the python object
def readPickle(fileName):
    infile = open(cwd + "/" + fileName + '.pickle','rb')
    data = pickle.load(infile)
    infile.close()
    return data


# adjustCoordinates(data)
# Changes coordinates in the square list so that the bottom left corner reflects the coordinates (0, 0)
# Returns: the data with the adjusted squares list
def adjustCoordinates(data):
    for index in range(len(data['squares'])):
        y1, y2, x1, x2 = data['squares'][index]
        data['squares'][index] = (
            x1,
            x2,
            data['height'] - y1 - 1,
            data['height'] - y2 - 1
        )
    return data


'''
######################## Obstacle Set Generator ######################## 
'''


# createVertex(XMLNode, int, int)
# Inserts a vertex XML node inside of the parent node with the coordinates x and y
# Returns: The created Vertex XML node
def createVertex(parent, x, y):
    node = ET.SubElement(parent, 'Vertex')
    node.set('p_x', str(x))
    node.set('p_y', str(y))
    return node


# createObstacle(XMLNode, tuple[4])
# Inserts a obstacle XML node inside of the parent node with 4 vertexes
# Returns: The created Obstacle XML node
def createObstacle(parent, square):
    node = ET.SubElement(parent, 'Obstacle')
    node.set('closed', '1')
    x1, x2, y1, y2 = square
    createVertex(node, x1, y2)
    createVertex(node, x2, y2)
    createVertex(node, x2, y1)
    createVertex(node, x1, y1)
    return node


# createBorder(XMLNode, int, int)
# Inserts a obstacle XML node inside of the parent node
# Note: Only use this for creating the obstacle for the outer ring of the room
# Returns: The created Obstacle XML node
def createBorder(parent, width, height):
    node = ET.SubElement(parent, 'Obstacle')
    node.set('closed', '1')
    createVertex(node, 0, height)
    createVertex(node, width, height)
    createVertex(node, width, 0)
    createVertex(node, 0, 0)
    return node


# createObstacleSet()
# Creates an ObstacleSet parent XML node filled with obstacles
# Returns: The created ObstacleSet node
def createObstacleSet(data):
    root = ET.Element('ObstacleSet')
    root.set('type', 'explicit')
    root.set('class', '1')
    createBorder(root, data['width'], data['height'])
    for square in data['squares']:
        createObstacle(root, square)
    return root


# writeToXML(XMLNode, string)
# Writes the given XML node to file
def writeToXML(node, fileName):
    data = ET.tostring(node)
    data = xdm.parseString(data)
    data = data.toprettyxml(indent="\t")
    outfile = open(cwd + "/" + fileName + '.xml', 'w')
    outfile.write(data)
    outfile.close()


'''
######################## Graph Generator ######################## 
'''


# createGraph(data)
# Generates a graph for all squares at least 4x4 in size
# Returns: a list of vertices and edges from the squares in data
def createGraph(data):
    v = []
    e = []
    for i in range(len(data['squares'])):
        x1, x2, y1, y2 = data['squares'][i]
        square_length = x2 - x1
        square_height = y1 - y2
        if square_length >= 4 and square_height >= 4:
            v += [
                (x1 + 1, y2 + 1),  # Top Left Vertex
                (x1 + 1, y1 - 1),  # Bottom Left Vertex
                (x2 - 1, y1 - 1),  # Top Right Vertex
                (x2 - 1, y2 + 1),  # Bottom Right Vertex
            ]
            i *= 4
            e += [
                (i, i + 1),
                (i, i + 2),
                (i, i + 3),
                (i + 1, i + 3),
                (i + 2, i + 1),
                (i + 2, i + 3),
            ]
    return v, e


# writeToTXT(int, int, string)
# Writes the list of vertices and edges to file
def writeToTXT(vertices, edges, fileName):
    total_vertices = len(vertices)
    total_edges = len(edges)
    outfile = open(cwd + "/" + fileName + '.txt', 'w')

    outfile.write(str(total_vertices) + "\n")
    for vertex in vertices:
        outfile.write("3 " + str(vertex[0]) + " " + str(vertex[1]) + "\n")

    outfile.write(str(total_edges) + "\n")
    for edge in edges:
        outfile.write(str(edge[0]) + " " + str(edge[1]) + "\n")

    outfile.close()


'''
######################## Main Program ######################## 
'''

# Main Program

file_data = readPickle(base_name)
file_data = adjustCoordinates(file_data)
print(file_data['width'])
print(file_data['height'])

if make_obstacles:
    obstacle_set = createObstacleSet(file_data)
    writeToXML(obstacle_set, base_name)

if make_graph:
    vertices, edges = createGraph(file_data)
    writeToTXT(base_name, vertices, edges)




