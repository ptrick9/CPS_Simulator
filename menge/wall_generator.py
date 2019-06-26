import os
import pickle
import xml.etree.ElementTree as ET
import xml.dom.minidom as xdm


'''
######################## Obstacle Set Generator ######################## 
'''


# create_vertex(XMLNode, int, int)
# Inserts a vertex XML node inside of the parent node with the coordinates x and y
# Returns: The created Vertex XML node
def create_vertex(parent, x, y):
    node = ET.SubElement(parent, 'Vertex')
    node.set('p_x', str(x))
    node.set('p_y', str(y))
    return node


# create_obstacle(XMLNode, tuple[4])
# Inserts a obstacle XML node inside of the parent node with 4 vertexes
# Returns: The created Obstacle XML node
def create_obstacle(parent, square):
    node = ET.SubElement(parent, 'Obstacle')
    node.set('closed', '1')
    x1, x2, y1, y2 = square
    create_vertex(node, x1, y2)
    create_vertex(node, x2, y2)
    create_vertex(node, x2, y1)
    create_vertex(node, x1, y1)
    return node


# create_border(XMLNode, int, int)
# Inserts a obstacle XML node inside of the parent node
# Note: Only use this for creating the obstacle for the outer ring of the room
# Returns: The created Obstacle XML node
def create_border(parent, width, height):
    node = ET.SubElement(parent, 'Obstacle')
    node.set('closed', '1')
    create_vertex(node, 0, height)
    create_vertex(node, width, height)
    create_vertex(node, width, 0)
    create_vertex(node, 0, 0)
    return node


# createObstacleSet(data)
# Creates an ObstacleSet parent XML node filled with obstacles
# Returns: The created ObstacleSet node
def create_obstacle_set(data):
    root = ET.Element('ObstacleSet')
    root.set('type', 'explicit')
    root.set('class', '1')
    create_border(root, data['width'], data['height'])
    for square in data['squares']:
        create_obstacle(root, square)
    return root


'''
######################## Main Program ######################## 
'''

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


# adjustCoordinates(data)
# Changes coordinates in the square list so that the bottom left corner reflects the coordinates (0, 0)
# Returns: the data with the adjusted squares list
def adjust_coordinates(data):
    for index in range(len(data['squares'])):
        y1, y2, x1, x2 = data['squares'][index]
        data['squares'][index] = (
            x1,
            x2,
            data['height'] - y1,
            data['height'] - y2
        )
    return data


# writeToXML(XMLNode, string)
# Writes the given XML node to file
def write_to_XML(node, fileName):
    data = ET.tostring(node)
    data = xdm.parseString(data)
    data = data.toprettyxml(indent="\t")
    outfile = open('%s/%s.xml' % (os.getcwd(), fileName), 'w')
    outfile.write(data)
    outfile.close()


def build(base_name, data):
    data = adjust_coordinates(data)
    obstacle_set = create_obstacle_set(data)
    write_to_XML(obstacle_set, base_name)
