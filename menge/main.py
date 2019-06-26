import scipy.misc
import copy
import random
import os


BASE_NAME = 'marathon_street_map'
OUTPUT_OBSTACLES = False
OUTPUT_GRAPH = True
OUTPUT_IMAGES = False

if OUTPUT_OBSTACLES:
    import wall_generator

if OUTPUT_GRAPH:
    import graph_generator

if OUTPUT_IMAGES:
    import imageio


# is_walkable(image, int, int)
# Returns: True if a pixel is white, else false
def is_walkable(image, x, y):
    if image[x][y][1] == 0:
        return True
    return False


# is_wall(image, int, int)
# Returns: True if a pixel is black, else false
def is_wall(image, x, y):
    if image[x][y][1] == 255:
        return True
    return False


# build_point_dict(image, function)
# image: image representation of the file
# func: a function that checks if a pixel is a certain color (is_wall or is_walkable)
# Returns: a dictionary which is a boolean grid representation of the image
def build_point_dict(image, func):
    point_dict = {}
    for x in range(image.shape[0]):
        for y in range(image.shape[1]):
            if not func(image, x, y):
                point_dict[(x, y)] = True
            else:
                point_dict[(x, y)] = False
    return point_dict


# remove_square(point_dict, x1, x2, y1, y2)
# Modifies the point_dict, removing all coordinates between the given points
def remove_square(point_dict, x1, x2, y1, y2):
    for i in range(x1, x2 + 1):
        for j in range(y1, y2 + 1):
            point_dict[(i, j)] = False


# build_square_list(image, point_dict)
# Returns: a list of squares where each square is a tuple[4] of coordinates representing top left and bottom
# right corners of a square
def build_square_list(image, point_dict):
    squares = []
    while True:

        top_left = None
        for k in sorted(point_dict.keys()):
            if point_dict[k]:
                top_left = list(k)
                break
        if top_left is None:
            break
        temp = copy.deepcopy(top_left)

        # [0] is x coord, [1] is y coord
        while temp[0] + 1 < image.shape[0] and point_dict[temp[0] + 1, temp[1]] == True:
            temp[0] += 1

        collide = False
        y_test = copy.deepcopy(top_left)

        while not collide:
            y_test[1] += 1
            if y_test[1] == image.shape[1]:
                collide = True
                break

            for x_val in range(top_left[0], temp[0]+1):
                if point_dict[x_val, y_test[1]] == False:
                    collide = True

        bottom_right = (temp[0], y_test[1] - 1)

        remove_square(point_dict, top_left[0], bottom_right[0], top_left[1], bottom_right[1])
        squares.append((top_left[0], bottom_right[0], top_left[1], bottom_right[1]))
    return squares


# build_border_dict(list(square))
# Returns: Dictionary containing graph of adjacent squares
def build_border_dict(squares):

    border_dict = {}
    for x in range(len(squares)):
        border_dict[x] = []

    for y in range(len(squares)):
        square = squares[y]

        for z in range(y + 1, len(squares)):
            new_square = squares[z]

            if new_square[0] >= square[0] and new_square[1] <= square[1]:
                if new_square[2] == square[3] + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

                elif new_square[3] == square[2] - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

            elif new_square[2] >= square[2] and new_square[3] <= square[3]:
                if new_square[0] == square[1] + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

                elif new_square[1] == square[0] - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

            if square[0] >= new_square[0] and square[1] <= new_square[1]:
                if square[2] == new_square[3] + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

                elif square[3] == new_square[2] - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

            elif square[2] >= new_square[2] and square[3] <= new_square[3]:
                if square[0] == new_square[1] + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

                elif square[1] == new_square[0] - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)

    return border_dict


# write_images(String, image, list(square))
# Colors in the square regions of an image and outputs them to file.
def write_images(name, image, squares):
    images = []
    for xx, square in enumerate(squares):

        for i in range(square[0], square[1] + 1):
            for j in range(square[2], square[3] + 1):
                image[i][j] = (color_list[xx])

        images.append(copy.deepcopy(image))
        scipy.misc.imsave('%s_%d.png' % (name, xx), image)

    imageio.mimsave('%s.gif' % name, images)


'''
######################## Main Program ######################## 
'''


im = scipy.misc.imread('%s/%s.png' % (os.getcwd(), BASE_NAME), mode='RGBA')
data = {
    'width': im.shape[1],
    'height': im.shape[0],
}

color_list = []
if OUTPUT_IMAGES:
    color_list = [(random.randrange(255), random.randrange(255), random.randrange(255), 128) for i in range(2000)]

if OUTPUT_OBSTACLES:
    print("Generating wall squares...")
    wall_points = build_point_dict(im, is_wall)
    wall_squares = build_square_list(im, wall_points)
    data['squares'] = wall_squares

    print("Writing walls to XML file...")
    wall_generator.build(BASE_NAME, data)

    if OUTPUT_IMAGES:
        print("Generating wall squares images...")
        write_images('%s_wall' % BASE_NAME, im, wall_squares)

if OUTPUT_GRAPH:
    print("Generating walkable squares...")
    walkable_points = build_point_dict(im, is_walkable)
    walkable_squares = build_square_list(im, walkable_points)
    data['squares'] = walkable_squares
    data['graph'] = build_border_dict(walkable_squares)

    print("Writing graph to TXT file...")
    graph_generator.build(BASE_NAME, data)

    if OUTPUT_IMAGES:
        print("Generating walkable squares images...")
        write_images('%s_walkable' % BASE_NAME, im, walkable_squares)

print("Done!")

