import scipy.misc
import copy
import random
import os
import pickle

localPC = True
text = False
if localPC:
    import imageio


#base_name = 'actual_circle_x4'
base_name = 'Basic'
#base_name = 'empty'

im = scipy.misc.imread(os.getcwd() + '/%s.png' % (base_name), mode='RGBA')

# print(im.shape)
# print(im[0][0])
# print(im[9][9])

color_list = [(random.randrange(255), random.randrange(255), random.randrange(255), 128) for i in range(2000)]
point_dict = {}
point_list = []
square_list = []
images = []


def isWall(im, x, y):
    if im[x][y][1] == 0:
        return True
    else:
        return False


def removeSquare(x1, x2, y1, y2):
    for i in range(x1, x2 + 1):
        for j in range(y1, y2 + 1):
            point_dict[(i, j)] = False
            #point_list.remove((i, j))

for x in range(im.shape[0]):
    for y in range(im.shape[1]):
        if not isWall(im, x, y):
            # point_list.append((x, y))
            point_dict[(x, y)] = True
        else:
            point_dict[(x, y)] = False

iter = 0

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

    # print(point_dict)
    while temp[0] + 1 < im.shape[0] and point_dict[temp[0] + 1, temp[1]] == True:
        temp[0] += 1

    print("found right")

    collide = False
    y_test = copy.deepcopy(top_left)

    while not collide:
        y_test[1] += 1
        if (y_test[1] == im.shape[1]):
            collide = True
            break

        for x_val in range(top_left[0], temp[0]+1):
            if point_dict[x_val, y_test[1]] == False:
                collide = True

    bottom_right = (temp[0], y_test[1] - 1)

    print((top_left[0], bottom_right[0], top_left[1], bottom_right[1]))
    removeSquare(top_left[0], bottom_right[0], top_left[1], bottom_right[1])
    square_list.append((top_left[0], bottom_right[0], top_left[1], bottom_right[1]))
    print((top_left[0], bottom_right[0], top_left[1], bottom_right[1]))

# print(square_list)





for xx,square in enumerate(square_list):

    for i in range(square[0], square[1] + 1):
        for j in range(square[2], square[3] + 1):
            im[i][j] = (color_list[xx])
    s = square
    #if localPC:
    #    bottomLeftCornerOfText = ((int(s[2]+(s[3]-s[2])/2), int(s[0]+(s[1]-s[0])/2)))

    images.append(copy.deepcopy(im))
    scipy.misc.imsave('%s_%d.png' % (base_name, iter), im)
    iter += 1

if localPC:
    imageio.mimsave('%s.gif' % base_name, images)

# print(point_list)
# print(square_list)
# print(len(square_list))

border_dict = {}
white_pixels = []

for x in range(len(square_list)):
    border_dict[x] = []

# squares go (x1, x2, y1, y2)
#             0   1   2   3
for y in range(len(square_list)):
    square = square_list[y]

    for z in range(y + 1, len(square_list)):
        new_square = square_list[z]

        if new_square[0] >= square[0] and new_square[1] <= square[1]:
            if new_square[2] == square[3] + 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((int((new_square[1] - new_square[0]) / 2) + new_square[0], new_square[2]))
                white_pixels.append((int((new_square[1] - new_square[0]) / 2) + new_square[0], square[3]))

            elif new_square[3] == square[2] - 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((int((new_square[1] - new_square[0]) / 2) + new_square[0], new_square[3]))
                white_pixels.append((int((new_square[1] - new_square[0]) / 2) + new_square[0], square[2]))

        elif new_square[2] >= square[2] and new_square[3] <= square[3]:
            if new_square[0] == square[1] + 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((new_square[0], int((new_square[3] - new_square[2]) / 2) + new_square[2]))
                white_pixels.append((square[1], int((new_square[3] - new_square[2]) / 2) + new_square[2]))

            elif new_square[1] == square[0] - 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((new_square[1], int((new_square[3] - new_square[2]) / 2) + new_square[2]))
                white_pixels.append((square[0], int((new_square[3] - new_square[2]) / 2) + new_square[2]))

        if square[0] >= new_square[0] and square[1] <= new_square[1]:
            if square[2] == new_square[3] + 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((int((square[1] - square[0]) / 2) + square[0], square[2]))
                white_pixels.append((int((square[1] - square[0]) / 2) + square[0], new_square[3]))

            elif square[3] == new_square[2] - 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((int((square[1] - square[0]) / 2) + square[0], square[3]))
                white_pixels.append((int((square[1] - square[0]) / 2) + square[0], new_square[2]))

        elif square[2] >= new_square[2] and square[3] <= new_square[3]:
            if square[0] == new_square[1] + 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((square[0], int((square[3] - square[2]) / 2) + square[2]))
                white_pixels.append((new_square[1], int((square[3] - square[2]) / 2) + square[2]))

            elif square[1] == new_square[0] - 1:
                border_dict[y].append(z)
                border_dict[z].append(y)

                white_pixels.append((square[1], int((square[3] - square[2]) / 2) + square[2]))
                white_pixels.append((new_square[0], int((square[3] - square[2]) / 2) + square[2]))

# print(border_dict)
# print(white_pixels)

pickle_out = {
    'width':    im.shape[1],
    'height':   im.shape[0],
    'squares':  square_list,
    'graph': border_dict
}

outfile = open(os.getcwd() + "/%s.pickle" % base_name, 'wb')
pickle.dump(pickle_out, outfile)
outfile.close()

x = {}




