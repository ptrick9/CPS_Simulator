
from math import sqrt
import pickle
import scipy.misc
import random
import os

color_list = [(random.randrange(255), random.randrange(255), random.randrange(255), 128) for i in range(2000)]
base_name = 'marathon_street_map'
im = scipy.misc.imread(os.getcwd() + '/%s.png' % (base_name), mode='RGBA')
image_num = 50

class Square:
    def __init__(self, x1, x2, y1, y2, can_cut, id_num):
        self.x1 = x1
        self.x2 = x2
        self.y1 = y1
        self.y2 = y2
        self.can_cut = can_cut
        self.id_num = id_num
        self.routers = []


    def __str__(self):
        return "id: %d x1: %d x2: %d y1: %d y2: %d can_cut: %s" % (
        self.id_num, self.x1, self.x2, self.y1, self.y2, self.can_cut)

    def __repr__(self):
        return "id: %d x1: %d x2: %d y1: %d y2: %d can_cut: %s" % (
        self.id_num, self.x1, self.x2, self.y1, self.y2, self.can_cut)

    def __eq__(self, other):
        if self.x1 == other.x1 and self.x2 == other.x2 and self.y1 == other.y1 and self.y2 == other.y2:
            return True
        else:
            return False

    def top_left(self):
        return self.y1, self.x1

    def top_right(self):
        return self.y2, self.x1

    def bottom_left(self):
        return self.y1, self.x2

    def bottom_right(self):
        return self.y2, self.x2


def side_ratio(sq1, sq2):
    if sq1.can_cut:
        ret = 1.1
        if sq1.x1 == sq2.x1 or sq1.x2 == sq2.x2:
            ret = (sq1.x2 - sq1.x1) / max((sq2.x2 - sq2.x1),1)
        elif sq1.y1 == sq2.y1 or sq1.y2 == sq2.y2:
            ret = (sq1.y2 - sq1.y1) / max((sq2.y2 - sq2.y1),1)
        if ret > 1:
            return 0
        else:
            return ret
    else:
        return 0


def area_ratio(sq1, sq2):
    if sq1.can_cut:
        if (sq1.x1 > sq2.x1 and sq1.x2 < sq2.x2):
            return (sq1.x2 - sq1.x1) / max((sq2.x2 - sq2.x1),1)
        elif (sq1.y1 > sq2.y1 and sq1.y2 < sq2.y2):
            return (sq1.y2 - sq1.y1) / max((sq2.y2 - sq2.y1),1)
        else:
            return 0
    else:
        return 0


def single_cut(sq1, sq2):
    new_squares = []

    if sq1.x1 == sq2.x1 and sq1.x2 == sq2.x2:
        if sq1.y1 < sq2.y1:
            new_squares.append(Square(sq1.x1, sq2.x2, sq1.y1, sq2.y2, True, sq1.id_num))
        else:
            new_squares.append(Square(sq1.x1, sq2.x2, sq2.y1, sq1.y2, True, sq1.id_num))

    elif sq1.y1 == sq2.y1 and sq1.y2 == sq2.y2:
        if sq1.x1 < sq2.x1:
            new_squares.append(Square(sq1.x1, sq2.x2, sq1.y1, sq1.y2, True, sq1.id_num))
        else:
            new_squares.append(Square(sq2.x1, sq1.x2, sq1.y1, sq1.y2, True, sq1.id_num))


    elif sq1.x1 == sq2.x1:  # cut on x2
        # print('x1', end=' ')
        if sq1.y1 < sq2.y1:
            new_squares.append(Square(sq1.x1, sq1.x2, sq1.y1, sq2.y2, True, sq1.id_num))
            new_squares.append(Square(sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, False, sq2.id_num))
            # print('we')
        else:
            new_squares.append(Square(sq1.x1, sq1.x2, sq2.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, False, sq2.id_num))
            # print('ew')

    elif sq1.x2 == sq2.x2:
        # print('x2', end=' ')
        if sq1.y1 < sq2.y1:
            new_squares.append(Square(sq1.x1, sq1.x2, sq1.y1, sq2.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, False, sq2.id_num))
            # print('we')
        else:
            new_squares.append(Square(sq1.x1, sq1.x2, sq2.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, False, sq2.id_num))
            # print('ew')

    elif sq1.y1 == sq2.y1:
        # print('y1', end=' ')
        if sq1.x1 < sq2.x1:
            new_squares.append(Square(sq1.x1, sq2.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, False, sq2.id_num))
            # print('ns')
        else:
            new_squares.append(Square(sq2.x1, sq1.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, False, sq2.id_num))
            # print('sn')

    elif sq1.y2 == sq2.y2:
        # print('y2', end=' ')
        if sq1.x1 < sq2.x1:
            new_squares.append(Square(sq1.x1, sq2.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, False, sq2.id_num))
            # print('ns')
        else:
            new_squares.append(Square(sq2.x1, sq1.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, False, sq2.id_num))
            # print('sn')

    print("CUT")
    return new_squares


def double_cut(sq1, sq2):
    new_squares = []
    if sq1.x1 > sq2.x1 and sq1.x2 < sq2.x2:
        if sq1.y2 + 1 == sq2.y1:
            new_squares.append(Square(sq1.x1, sq1.x2, sq1.y1, sq2.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, False, sq2.id_num))
            new_squares.append(Square(sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, False, -1))
        else:
            new_squares.append(Square(sq1.x1, sq1.x2, sq2.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, False, sq2.id_num))
            new_squares.append(Square(sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, False, -1))


    elif sq1.y1 > sq2.y1 and sq1.y2 < sq2.y2:
        if sq1.x2 + 1 == sq2.x1:
            new_squares.append(Square(sq1.x1, sq2.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, False, sq2.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, False, -1))
        else:
            new_squares.append(Square(sq2.x1, sq1.x2, sq1.y1, sq1.y2, True, sq1.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, False, sq2.id_num))
            new_squares.append(Square(sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, False, -1))

    return new_squares


def rebuild(sq_list):
    for x in range(len(square_list)):
        border_dict[x] = []
        square_list[x].routers = []

    for y in range(len(sq_list)):
        square = sq_list[y]

        for z in range(y + 1, len(sq_list)):
            new_square = sq_list[z]

            if new_square.x1 >= square.x1 and new_square.x2 <= square.x2:
                if new_square.y1 == square.y2 + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        new_square.y1,
                        int((new_square.x2 - new_square.x1) / 2) + new_square.x1,
                    ))
                    square_list[y].routers.append((
                        square.y2,
                        int((new_square.x2 - new_square.x1) / 2) + new_square.x1,
                    ))

                elif new_square.y2 == square.y1 - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        new_square.y2,
                        int((new_square.x2 - new_square.x1) / 2) + new_square.x1,
                    ))
                    square_list[y].routers.append((
                        square.y1,
                        int((new_square.x2 - new_square.x1) / 2) + new_square.x1,
                    ))

            elif new_square.y1 >= square.y1 and new_square.y2 <= square.y2:
                if new_square.x1 == square.x2 + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        int((new_square.y2 - new_square.y1) / 2) + new_square.y1,
                        new_square.x1,
                    ))
                    square_list[y].routers.append((
                        int((new_square.y2 - new_square.y1) / 2) + new_square.y1,
                        square.x2,

                    ))

                elif new_square.x2 == square.x1 - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        int((new_square.y2 - new_square.y1) / 2) + new_square.y1,
                        new_square.x2,
                    ))
                    square_list[y].routers.append((
                        int((new_square.y2 - new_square.y1) / 2) + new_square.y1,
                        square.x1,
                    ))

            if square.x1 >= new_square.x1 and square.x2 <= new_square.x2:
                if square.y1 == new_square.y2 + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        square.y1,
                        int((square.x2 - square.x1) / 2) + square.x1,
                    ))
                    square_list[y].routers.append((
                        new_square.y2,
                        int((square.x2 - square.x1) / 2) + square.x1,
                    ))

                elif square.y2 == new_square.y1 - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        square.y2,
                        int((square.x2 - square.x1) / 2) + square.x1,
                    ))
                    square_list[y].routers.append((
                        new_square.y1,
                        int((square.x2 - square.x1) / 2) + square.x1,
                    ))

            elif square.y1 >= new_square.y1 and square.y2 <= new_square.y2:
                if square.x1 == new_square.x2 + 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        int((square.y2 - square.y1) / 2) + square.y1,
                        square.x1,
                    ))
                    square_list[y].routers.append((
                        int((square.y2 - square.y1) / 2) + square.y1,
                        new_square.x2,

                    ))

                elif square.x2 == new_square.x1 - 1:
                    border_dict[y].append(z)
                    border_dict[z].append(y)
                    square_list[z].routers.append((
                        int((square.y2 - square.y1) / 2) + square.y1,
                        square.x2,
                    ))
                    square_list[y].routers.append((
                        int((square.y2 - square.y1) / 2) + square.y1,
                        new_square.x1,
                   ))
    global image_num
    for xx, sq in enumerate(square_list):

        for i in range(sq.x1, sq.x2 + 1):
            for j in range(sq.y1, sq.y2 + 1):
                im[i][j] = (color_list[xx])
        # if localPC:
        #    bottomLeftCornerOfText = ((int(s[2]+(s[3]-s[2])/2), int(s[0]+(s[1]-s[0])/2)))

    scipy.misc.imsave('%s_%d.png' % (base_name, image_num), im)
    image_num += 1


f = open(base_name + '.pickle', 'rb')
data = pickle.load(f)

coord_list = data['squares']
border_dict = data['graph']

square_list = []

for i, s in enumerate(coord_list):
    square_list.append(Square(s[0], s[1], s[2], s[3], True, i))

# print(square_list[5])
# print(square_list[12])
# print(square_list[6])

# rebuild(square_list)

while True:
    i = 0
    rebuilt = False

    while i < len(square_list) and not rebuilt:

        for n in border_dict[i]:

            s_rat = side_ratio(square_list[i], square_list[n])

            if s_rat > 0.6:
                new_squares = single_cut(square_list[i], square_list[n])

                s1 = square_list[n]
                s2 = square_list[i]

                square_list.remove(s1)
                square_list.remove(s2)

                square_list.extend(new_squares)

                rebuild(square_list)

                rebuilt = True



                break

            a_rat = area_ratio(square_list[i], square_list[n])

            if a_rat > 0.6:
                new_squares = double_cut(square_list[i], square_list[n])

                s1 = square_list[n]
                s2 = square_list[i]

                square_list.remove(s1)
                square_list.remove(s2)

                square_list.append(new_squares[0])
                square_list.append(new_squares[1])

                new_squares[2].id_num = len(square_list)

                square_list.append(new_squares[2])

                rebuild(square_list)

                rebuilt = True

                break

        i += 1

    if not rebuilt:
        break;

rebuild(square_list)
print("\nPrinting Square List and Border Dict")
print(square_list)
print(border_dict)
print()


simulation_width = data['width']
simulation_height = data['height']
vertices = []
edges = []
v_sq_map = {}
vmap = {}


# Checks if the two coordinate pairs are nearby
def is_adjacent(c1, c2):
    x1, y1 = c1
    x2, y2 = c2
    if sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1)) < 2:
        return True
    return False


# Adding corner vertices to each square and connecting them to themselves
for i, s in enumerate(square_list):

    s.x1 = simulation_height - s.x1 - 1
    s.x2 = simulation_height - s.x2 - 1
    degree = 3 + len(s.routers)
    vertices += [
        (degree, s.y1, s.x2),  # Top Left Vertex
        (degree, s.y1, s.x1),  # Bottom Left Vertex
        (degree, s.y2, s.x1),  # Top Right Vertex
        (degree, s.y2, s.x2),  # Bottom Right Vertex
    ]



    i *= 4
    edges += [
        (i, i + 1),
        (i, i + 2),
        (i, i + 3),
        (i + 1, i + 2),
        (i + 1, i + 3),
        (i + 2, i + 3),
    ]

# Adding router vertices to each of the squares and connecting them to the corner vertices
for i, s in enumerate(square_list):
    index = i * 4
    v_sq_map[i] = []
    for ri, coord in enumerate(s.routers):
        x, y = coord
        y = simulation_height - y - 1
        vid = len(vertices)
        degree = 5 + len(s.routers) - 1
        vertices += [
            (degree, x, y)
        ]
        edges += [
            (vid, index),
            (vid, index + 1),
            (vid, index + 2),
            (vid, index + 3),
        ]
        vmap[(x, y)] = vid
        v_sq_map[i].append(vid)

        # Connecting internal routers
        if len(s.routers) >= 2:
            for last_rid in range(ri - 1, -1, -1):
                x, y = s.routers[last_rid]
                y = simulation_height - y - 1
                edges += [
                    (vmap[(x, y)], vid)
                ]



# Connecting adjacent routers
for i, s in enumerate(square_list):

    for x, y in s.routers:
        y = simulation_height - y - 1

        for adjacent_sid in border_dict[i]:

            for vid in v_sq_map[adjacent_sid]:

                c1 = (x, y)
                c2 = (vertices[vid][1], vertices[vid][2])
                if is_adjacent(c1, c2) and (vmap[c2], vmap[c1]) not in edges:
                    edges += [
                        (vmap[c1], vmap[c2])
                    ]


total_vertices = len(vertices)
total_edges = len(edges)
outfile = open(base_name + '.txt', 'w')

outfile.write(str(total_vertices) + "\n")
for vertex in vertices:
    outfile.write("%s %s %s\n" % (vertex[0], vertex[1], vertex[2]))

outfile.write(str(total_edges) + "\n")
for edge in edges:
    outfile.write(str(edge[0]) + " " + str(edge[1]) + "\n")

outfile.close()




