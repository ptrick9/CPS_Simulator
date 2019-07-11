import scipy.misc
import copy
import random
import imageio


# is_wall(image, int, int)
# Returns: True if a pixel is black, else false
def is_wall(image, color, x, y):
    if image[x][y][0] == color \
            and image[x][y][1] == color \
            and image[x][y][2] == color:
        return True
    else:
        return False


# build_point_dict(image, function)
# image: image representation of the file
# func: a function that checks if a pixel is a certain color (is_wall or is_walkable)
# Returns: a dictionary which is a boolean grid representation of the image
def build_point_dict(image, color):
    point_dict = {}
    for x in range(image.shape[0]):
        for y in range(image.shape[1]):
            if not is_wall(image, color, x, y):
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
        while temp[1] + 1 < image.shape[1] and point_dict[temp[0], temp[1] + 1] == True:
            temp[1] += 1

        collide = False
        y_test = copy.deepcopy(top_left)
        while not collide:
            y_test[0] += 1
            if y_test[0] >= image.shape[0]:
                collide = True
                break
            for x_val in range(top_left[1], temp[1] + 1):
                if point_dict[y_test[0], x_val] == False:
                    collide = True

        bottom_right = (temp[1], y_test[0] - 1)
        remove_square(point_dict, top_left[0], bottom_right[1], top_left[1], bottom_right[0])
        squares.append((top_left[0], bottom_right[1], top_left[1], bottom_right[0]))
    return squares


# build_border_dict(list(square))
# Returns: Dictionary containing graph of adjacent squares
def build_border_dict(squares):

    border_dict = {}
    for i in range(len(squares)):
        border_dict[i] = []

    return border_dict


# write_images(String, image, list(square))
# Colors in the square regions of an image and outputs them to file.
def write_images(name, image, squares):
    color_list = [(random.randrange(255), random.randrange(255), random.randrange(255), 128) for i in range(2000)]
    images = []
    for xx, square in enumerate(squares):
        for i in range(square[0], square[1] + 1):
            for j in range(square[2], square[3] + 1):
                image[i][j] = (color_list[xx])

        images.append(copy.deepcopy(image))
        scipy.misc.imsave('%s_%d.png' % (name, xx), image)

    imageio.mimsave('%s.gif' % name, images)
