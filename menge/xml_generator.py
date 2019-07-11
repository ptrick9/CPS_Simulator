import xml.etree.ElementTree as et

'''
######################## Behavior XML ######################## 
'''


''' ----------- Goals ----------- '''


def make_goal_set(id, capacity, locations):
    node = et.Element('GoalSet')
    node.set('id', str(id))
    for i, location in enumerate(locations):
        goal = et.SubElement(node, 'Goal')
        goal.set('type', 'circle')
        goal.set('id', str(i))
        goal.set('capacity', str(capacity))
        goal.set('x', str(location[0]))
        goal.set('y', str(location[1]))
        goal.set('radius', '1')
    return node


''' ----------- States ----------- '''


def make_state_static(name):
    node = et.Element('State')
    node.set('name', name)
    node.set('final', '0')
    goal_selector = et.SubElement(node, 'GoalSelector')
    goal_selector.set('type', 'identity')
    vel_component = et.SubElement(node, 'VelComponent')
    vel_component.set('type', 'zero')
    return node


def make_state_teleport(name, location):
    node = make_state_static(name)
    action = et.SubElement(node, 'Action')
    action.set('type', 'teleport')
    action.set('dist', 'c')
    action.set('x_value', str(location[0]))
    action.set('y_value', str(location[1]))
    return node


def make_state_travel(name, goal_set_id, map_name):
    node = et.Element('State')
    node.set('name', name)
    node.set('final', '0')
    goal_selector = et.SubElement(node, 'GoalSelector')
    goal_selector.set('type', 'random')
    goal_selector.set('goal_set', str(goal_set_id))
    vel_component = et.SubElement(node, 'VelComponent')
    vel_component.set('type', 'road_map')
    vel_component.set('file_name', '%s.txt' % map_name)
    return node


''' ----------- Transitions  ----------- '''


def make_transition_timer(fro, to, min, max):
    node = et.Element('Transition')
    node.set('from', fro)
    node.set('to', to)
    condition = et.SubElement(node, 'Condition')
    condition.set('type', 'timer')
    condition.set('dist', 'u')
    condition.set('min', str(min))
    condition.set('max', str(max))
    condition.set('per_agent', '1')
    return node


def make_transition_goal_reached(fro, to):
    node = et.Element('Transition')
    node.set('from', fro)
    node.set('to', to)
    condition = et.SubElement(node, 'Condition')
    condition.set('type', 'goal_reached')
    condition.set('distance', '1')
    return node


def make_transition_random(fro, destinations):
    node = et.Element('Transition')
    node.set('from', fro)
    condition = et.SubElement(node, 'Condition')
    condition.set('type', 'auto')
    target = et.SubElement(node, 'Target')
    target.set('type', 'prob')
    for destination in destinations:
        state = et.SubElement(target, 'State')
        state.set('name', destination[0])
        state.set('weight', str(destination[1]))
    return node


'''
######################## Scene XML ######################## 
'''


def make_agent_profile(id, speed):
    node = et.Element('AgentProfile')
    node.set('name', 'group_%d' % id)
    node.set('inherits', 'base')
    properties = et.SubElement(node, 'Common')
    properties.set('class', str(id))
    properties.set('pref_speed', str(speed))
    return node


def make_agent_group(id, amount):
    node = et.Element('AgentGroup')
    profile_selector = et.SubElement(node, 'ProfileSelector')
    profile_selector.set('type', 'const')
    profile_selector.set('name', 'group_%d' % id)
    state_selector = et.SubElement(node, 'StateSelector')
    state_selector.set('type', 'const')
    state_selector.set('name', 'Start_%d' % id)
    generator = et.SubElement(node, 'Generator')
    generator.set('type', 'rect_grid')
    generator.set('anchor_x', '-10')
    generator.set('anchor_y', '-10')
    generator.set('offset_x', '1')
    generator.set('offset_y', '0')
    generator.set('count_x', str(amount))
    generator.set('count_y', '1')
    return node

