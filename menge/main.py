import os
import imageio
import scipy.misc
import sys
import xml.dom.minidom as xdm
import xml.etree.ElementTree as et
import graph_generator
import square_generator
import wall_generator
import xml_generator

# Outputs to a folder titled MAIN_FILE_NAME

MAIN_FILE_NAME = 'stadium_seating'
WALL_FILE_NAME = 'stadium_seating_noPoints'

OUTPUT_SCENE_XML = True
OUTPUT_BEHAVIOR_XML = True
OUTPUT_VIEWER_XML = True
OUTPUT_LINK_XML = True
OUTPUT_GRAPH_TXT = True

MAIN_IMAGE = imageio.imread('%s.png' % MAIN_FILE_NAME)
WALL_IMAGE = scipy.misc.imread('%s.png' % WALL_FILE_NAME, mode='RGBA')

MAIN_XML = et.parse('%s.xml' % MAIN_FILE_NAME).getroot()
SCENE_XML = et.parse('base_scene.xml').getroot()
BEHAVIOR_XML = et.Element('BFSM')
VIEWER_XML = et.parse('base_viewer.xml').getroot()
TOTAL_GOAL_SETS = 0
COLOR_DICTIONARY = {}

''' === Main === '''

# create_color_dictionary(Image):
# image: an imageio image representing a scenario
# Returns: a dictionary where the key is rgb and the value is a list of coordinates the color can be found
def create_color_dictionary(image):
    height = image.shape[0]
    width = image.shape[1]
    color_coords = {}
    for x in range(width):
        for y in range(height):
            rgb = tuple(image[y][x][:3])
            if color_coords.get(rgb) is None:
                color_coords[rgb] = []
            color_coords[rgb].append((x, height - y))
    return color_coords


def get_node_rgb(node):
    r = int(node.attrib['r'])
    g = int(node.attrib['g'])
    b = int(node.attrib['b'])
    return tuple((r, g, b))


def add_agents(group_id, group_node):
    try:
        speed = int(group_node.attrib['speed'])
        node = xml_generator.make_agent_profile(group_id, speed)
        SCENE_XML.append(node)
    except KeyError:
        print("[ERROR] Failed to create agent profile for group %d. Missing 'speed' attribute." % group_id)
        return False
    except ValueError:
        print("[ERROR] Failed to create agent profile for group %d. Invalid type for 'speed' attribute." % group_id)
        return False

    try:
        amount = int(group_node.attrib['amount'])
        node = xml_generator.make_agent_group(group_id, amount)
        SCENE_XML.append(node)
    except KeyError:
        print("[ERROR] Failed to create agent group for group %d. Missing 'amount' attribute." % group_id)
        return False
    except ValueError:
        print("[ERROR] Failed to create agent group for group %d. Invalid type for 'amount' attribute." % group_id)
        return False

    return True


def add_goal_sets(group_id, group_node):

    for goal_set_id, goal_set_node in enumerate(group_node.findall('GoalSet')):
        capacity = None
        rgb = None
        destinations = None

        try:
            capacity = int(goal_set_node.attrib['capacity'])
        except KeyError:
            print("[ERROR] Failed to create goal set %d for group %d. Missing 'capacity' attribute." %
                  (goal_set_id, group_id))
            return False
        except ValueError:
            print("[ERROR] Failed to create goals et %d for group %d. Invalid type for 'capacity' attribute." %
                  (goal_set_id, group_id))
            return False

        try:
            rgb = get_node_rgb(goal_set_node.find('Color'))
        except KeyError:
            print("[ERROR] Failed to create goal set %d for group %d. Missing 'r', 'g', 'b' attributes." %
                  (goal_set_id, group_id))
            return False
        except ValueError:
            print("[ERROR] Failed to create goal set %d for group %d. Invalid type for 'r', 'g', 'b' attributes." %
                  (goal_set_id, group_id))
            return False

        try:
            destinations = COLOR_DICTIONARY[rgb]
        except KeyError:
            print("[ERROR] Failed to create goal set %d for group %d. Could not find pixels with an RGB value of %s in %s.png." %
                  (goal_set_id, group_id, rgb, MAIN_FILE_NAME))
            return False

        node = xml_generator.make_goal_set(goal_set_id + TOTAL_GOAL_SETS, capacity, destinations)
        BEHAVIOR_XML.append(node)
    return True


def add_spawn(group_id, group_node):

    min = None
    max = None
    rgb = None
    locations = []
    destinations = []

    try:
        min = int(group_node.find('Spawn').attrib['min'])
    except KeyError:
        print("[ERROR] Failed to create spawn for group %d. Missing 'min' attribute." % group_id)
        return False
    except ValueError:
        print("[ERROR] Failed to create spawn for group %d. Invalid type for 'min' attribute." % group_id)
        return False

    try:
        max = int(group_node.find('Spawn').attrib['max'])
    except KeyError:
        print("[ERROR] Failed to create spawn for group %d. Missing 'max' attribute." % group_id)
        return False

    try:
        rgb = get_node_rgb(group_node.find('Spawn').find('Color'))
    except KeyError:
        print("[ERROR] Failed to create spawn for group %d. Missing 'r', 'g', 'b' attributes." % group_id)
        return False
    except ValueError:
        print("[ERROR] Failed to create spawn for group %d. Invalid type for 'r', 'g', 'b' attributes." % group_id)
        return False

    state_start_name = 'Start_%d' % group_id
    start_node = xml_generator.make_state_static(state_start_name)
    BEHAVIOR_XML.append(start_node)

    state_wait_name = 'Start_Wait_%d' % group_id
    wait_node = xml_generator.make_state_static(state_wait_name)
    BEHAVIOR_XML.append(wait_node)

    start_wait_tran = xml_generator.make_transition_timer(state_start_name, state_wait_name, min, max)
    BEHAVIOR_XML.append(start_wait_tran)

    try:
        for spawn_id, location in enumerate(COLOR_DICTIONARY[rgb]):
            name = 'Spawn_%d_%d' % (group_id, spawn_id)
            locations.append(tuple((name, '1')))
            spawn_node = xml_generator.make_state_teleport(name, location)
            BEHAVIOR_XML.append(spawn_node)
    except KeyError:
        print("[ERROR] Failed to create spawn for group %d. Could not find pixels with a RGB value of %s in %s." %
              (group_id, rgb, MAIN_FILE_NAME))
        return False

    spawn_tran = xml_generator.make_transition_random(state_wait_name, locations)
    BEHAVIOR_XML.append(spawn_tran)

    for destination_id, destination_node in enumerate(group_node.find('Spawn').findall('Transition')):
        to = None
        chance = None

        try:
            to = int(destination_node.attrib['to']) + TOTAL_GOAL_SETS
        except KeyError:
            print("[ERROR] Failed to create spawn for group %d. Missing 'to' attribute for transition %d." %
                  (group_id, destination_id))
            return False
        except ValueError:
            print(
                "[ERROR] Failed to create spawn for group %d. Invalid type for 'to' attribute for transition %d." %
                (group_id, destination_id))
            return False

        try:
            chance = float(destination_node.attrib['chance'])
            destinations.append(tuple(('Travel_%d_%d' % (group_id, to), chance)))
        except KeyError:
            print("[ERROR] Failed to create spawn for group %d. Missing 'chance' attribute for transition %d." %
                  (group_id, destination_id))
            return False
        except ValueError:
            print(
                "[ERROR] Failed to create spawn for group %d. Invalid type for 'chance' attribute for transition %d." %
                (group_id, destination_id))
            return False

        spawn_names = [x[0] for x in locations]
        tran = xml_generator.make_transition_random(','.join(spawn_names), destinations)
        BEHAVIOR_XML.append(tran)
    return True


def add_goals(group_id, group_node):

    for i, goal_set_node in enumerate(group_node.findall('GoalSet')):
        goal_set_id = i + TOTAL_GOAL_SETS

        state_travel_name = 'Travel_%d_%d' % (group_id, goal_set_id)
        travel_node = xml_generator.make_state_travel(state_travel_name, goal_set_id, MAIN_FILE_NAME)
        BEHAVIOR_XML.append(travel_node)

        state_arrive_name = 'Arrive_%d_%d' % (group_id, goal_set_id)
        arrive_node = xml_generator.make_state_static(state_arrive_name)
        BEHAVIOR_XML.append(arrive_node)

        travel_arrive_tran = xml_generator.make_transition_goal_reached(state_travel_name, state_arrive_name)
        BEHAVIOR_XML.append(travel_arrive_tran)

        # Wait at goal state
        state_wait_name = 'Wait_%d_%d' % (group_id, goal_set_id)
        wait_node = xml_generator.make_state_static(state_wait_name)
        BEHAVIOR_XML.append(wait_node)

        min = None
        max = None
        try:
            min = int(goal_set_node.attrib['min'])
        except KeyError:
            print("[ERROR] Failed to create behavior for group %d. Missing 'min' attribute for goal set %d." %
                  (group_id, goal_set_id))
            return False
        except ValueError:
            print(
                "[ERROR] Failed to create behavior for group %d. Invalid type for 'min' attribute for goal set %d." %
                group_id, goal_set_id)
            return False

        try:
            max = int(goal_set_node.attrib['max'])
        except KeyError:
            print("[ERROR] Failed to create behavior for group %d. Missing 'max' attribute for goal set %d." %
                  (group_id, goal_set_id))
            return False
        except ValueError:
            print(
                "[ERROR] Failed to create behavior for group %d. Invalid type for 'max' attribute for goal set %d." %
                group_id, goal_set_id)
            return False

        timer_tran = xml_generator.make_transition_timer(state_arrive_name, state_wait_name, min, max)
        BEHAVIOR_XML.append(timer_tran)

        next_destinations = []
        for transition_id, transition_node in enumerate(goal_set_node.findall('Transition')):
            to = None
            chance = None
            try:
                to = int(transition_node.attrib['to'])
            except KeyError:
                print("[ERROR] Failed to add transition %d for goal set %d in group %d. Missing 'to' attribute." %
                      (transition_id, goal_set_id, group_id))
                return False
            except ValueError:
                print("[ERROR] Failed to add transition %d for goal set %d in group %d. Invalid type for 'to' attribute." %
                      (transition_id, goal_set_id, group_id))
                return False

            try:
                chance = float(transition_node.attrib['chance'])
            except KeyError:
                print("[ERROR] Failed to add transition %d for goal set %d in group %d. Missing 'chance' attribute." %
                      (transition_id, goal_set_id, group_id))
                return False
            except ValueError:
                print("[ERROR] Failed to add transition %d for goal set %d in group %d. Invalid type for 'chance' attribute." %
                    (transition_id, goal_set_id, group_id))
                return False

            if goal_set_id == int(to):
                to = to + TOTAL_GOAL_SETS
                next_destinations.append(tuple(('Arrive_%d_%d' % (group_id, to), chance)))
            else:
                to = to + TOTAL_GOAL_SETS
                next_destinations.append(tuple(('Travel_%d_%d' % (group_id, to), chance)))

        destination_tran = xml_generator.make_transition_random(state_wait_name, next_destinations)
        BEHAVIOR_XML.append(destination_tran)


    return True


def create_XML_link():
    root = et.Element('Project')
    root.set('scene', '%sS.xml' % MAIN_FILE_NAME)
    root.set('behavior', '%sB.xml' % MAIN_FILE_NAME)
    root.set('view', '%sV.xml' % MAIN_FILE_NAME)
    root.set('model', 'orca')
    root.set('dumpPath', 'images/%s' % MAIN_FILE_NAME)
    return root


def create_XML_scene_behavior():
    global TOTAL_GOAL_SETS
    for group_id, group_node in enumerate(MAIN_XML.findall('Group')):
        if OUTPUT_SCENE_XML:
            if not add_agents(group_id, group_node):
                break

        if OUTPUT_BEHAVIOR_XML:
            if not add_goal_sets(group_id, group_node):
                break

            if not add_spawn(group_id, group_node):
                break

            if not add_goals(group_id, group_node):
                break

            TOTAL_GOAL_SETS += len(group_node.findall('GoalSet'))

def create_XML_viewer():

    x = MAIN_IMAGE.shape[1] / 2 - 50
    y = MAIN_IMAGE.shape[0] / 2 - 50
    xtgt = x
    ytgt = y + .01
    scale = .5

    camera = VIEWER_XML.find('Camera')
    camera.set('xpos', str(x))
    camera.set('ypos', str(y))
    camera.set('xtgt', str(xtgt))
    camera.set('ytgt', str(ytgt))
    camera.set('orthoScale', str(scale))


def write_to_XML(node, fileName):
    data = et.tostring(node)
    data = xdm.parseString(data)
    data = data.toprettyxml(indent="\t")
    outfile = open('%s.xml' % fileName, 'w')
    outfile.write(data)
    outfile.close()


def print_progress(title, progress):
    length = 20
    block = int(round(length * progress))
    msg = "\r{0}: [{1}] {2}%".format(title, "#" * block + "-" * (length - block), round(progress * 100, 2))
    if progress >=1:
        msg += " Done\r\n"
    sys.stdout.write(msg)
    sys.stdout.flush()


'''
######################## Main Program ######################## 
'''

if OUTPUT_BEHAVIOR_XML or OUTPUT_SCENE_XML or OUTPUT_VIEWER_XML:
    print("Creating output directory '%s/'" % MAIN_FILE_NAME)
    path = '%s/' % MAIN_FILE_NAME
    if not os.path.exists(path):
        os.makedirs(path)

if OUTPUT_LINK_XML:
    print("Creating file '%s/%s.xml'..." % (MAIN_FILE_NAME, MAIN_FILE_NAME))
    link_xml = create_XML_link()
    write_to_XML(link_xml, '%s/%s' % (MAIN_FILE_NAME, MAIN_FILE_NAME))

if OUTPUT_VIEWER_XML:
    create_XML_viewer()
    print("Creating file '%s/%sV.xml'...." % (MAIN_FILE_NAME, MAIN_FILE_NAME))
    write_to_XML(VIEWER_XML, '%s/%sV' % (MAIN_FILE_NAME, MAIN_FILE_NAME))

if OUTPUT_BEHAVIOR_XML or OUTPUT_SCENE_XML:
    COLOR_DICTIONARY = create_color_dictionary(MAIN_IMAGE)
    create_XML_scene_behavior()

if OUTPUT_BEHAVIOR_XML:
    print("Creating file '%s/%sB.xml'..." % (MAIN_FILE_NAME, MAIN_FILE_NAME))
    write_to_XML(BEHAVIOR_XML, '%s/%sB' % (MAIN_FILE_NAME, MAIN_FILE_NAME))

if OUTPUT_SCENE_XML:
    print("Generating walls...")
    wall_points = square_generator.build_point_dict(WALL_IMAGE, 255)
    wall_squares = square_generator.build_square_list(WALL_IMAGE, wall_points)
    data = {
        'width': WALL_IMAGE.shape[1],
        'height': WALL_IMAGE.shape[0],
        'squares': wall_squares
    }

    obstacle_set_node = wall_generator.create_obstacle_set(data)
    SCENE_XML.append(obstacle_set_node)
    print("Creating file '%s/%sS.xml'..." % (MAIN_FILE_NAME, MAIN_FILE_NAME))
    write_to_XML(SCENE_XML, '%s/%sS' % (MAIN_FILE_NAME, MAIN_FILE_NAME))

if OUTPUT_GRAPH_TXT:
    print("Generating graph...")
    walkable_points = square_generator.build_point_dict(WALL_IMAGE, 0)
    walkable_squares = square_generator.build_square_list(WALL_IMAGE, walkable_points)
    data = {
        'width': WALL_IMAGE.shape[1],
        'height': WALL_IMAGE.shape[0],
        'squares': walkable_squares,
        'graph': square_generator.build_border_dict(walkable_squares),
    }

    print("Creating file %s/%s.txt" % (MAIN_FILE_NAME, MAIN_FILE_NAME))
    graph_generator.build("%s/%s" % (MAIN_FILE_NAME,MAIN_FILE_NAME), data)

print("Done!")
