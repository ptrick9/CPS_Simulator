from math import sqrt


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

    #print("CUT")
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


def rebuild(square_list, sq_list, border_dict):
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


# build_squares(list(tuple[4])
# old_squares are tuple[4] coordinates of top left and bottom right squares
# Returns: a list of square objects
def build_squares(old_squares):
    square_list = []
    for i, s in enumerate(old_squares):
        square_list.append(Square(s[0], s[1], s[2], s[3], True, i))
    return square_list


# merge_squares(list(square), border_dict)
# Reduces the amount of squares by merging overlapping regions
def merge_squares(square_list, border_dict):
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
                    rebuild(square_list, square_list, border_dict)
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
                    rebuild(square_list, square_list, border_dict)
                    rebuilt = True
                    break

            i += 1

        if not rebuilt:
            break;

    rebuild(square_list, square_list, border_dict)


# Checks if the two coordinate pairs are nearby
def is_adjacent(c1, c2):
    x1, y1 = c1
    x2, y2 = c2
    if sqrt((x2 - x1) * (x2 - x1) + (y2 - y1) * (y2 - y1)) < 2:
        return True
    return False


# add_internal_squares(list(square), list(vertex), list(edge)
# Adds vertices to the corner of each square and connects them with edges
# Modifies vertices and edges
def add_internal_squares(squares, vertices, edges, sim_height):
    for i, s in enumerate(squares):
        s.x1 = sim_height - s.x1
        s.x2 = sim_height - s.x2
        degree = 3 + len(s.routers)
        vertices.extend([
            (degree, s.y1, s.x2),  # Top Left Vertex
            (degree, s.y1, s.x1),  # Bottom Left Vertex
            (degree, s.y2, s.x1),  # Top Right Vertex
            (degree, s.y2, s.x2),  # Bottom Right Vertex
        ])

        i *= 4
        edges.extend([
            (i, i + 1),
            (i, i + 2),
            (i, i + 3),
            (i + 1, i + 2),
            (i + 1, i + 3),
            (i + 2, i + 3),
        ])


# add_routers(list(square), list(vertex), list(edge), dict)
# Adding router vertices to each of the squares and connecting them to the corner vertices
# Modifies vertices and edges adding the squares routers
# Modifies vmap by storing the vertex id of the newly added routers
# Modifies v_sq_map by storing a list of vertices ids for each square
def add_routers(squares, vertices, edges, vmap, v_sq_map, sim_height):
    for i, s in enumerate(squares):
        index = i * 4
        v_sq_map[i] = []
        for ri, coord in enumerate(s.routers):
            x, y = coord
            y = sim_height - y
            vid = len(vertices)
            degree = 5 + len(s.routers) - 1
            vertices.extend([
                (degree, x, y)
            ])
            edges.extend([
                (vid, index),
                (vid, index + 1),
                (vid, index + 2),
                (vid, index + 3),
            ])
            vmap[(x, y)] = vid
            v_sq_map[i].append(vid)

            # Connecting internal routers
            if len(s.routers) >= 2:
                for last_rid in range(ri - 1, -1, -1):
                    x, y = s.routers[last_rid]
                    y = sim_height - y
                    edges.extend([
                        (vmap[(x, y)], vid)
                    ])


# connect_adjacent_routers(list(square), list(vertex), list(edge), vmap, v_sq_map, int, border_dict)
# Connecting adjacent routers
# Modifies vertices, edges by adding connections to adjacent routers
def connect_adjacent_routers(squares, vertices, edges, vmap, v_sq_map, sim_height, border_dict):
    for i, s in enumerate(squares):
        for x, y in s.routers:
            y = sim_height - y
            for adjacent_sid in border_dict[i]:
                for vid in v_sq_map[adjacent_sid]:
                    c1 = (x, y)
                    c2 = (vertices[vid][1], vertices[vid][2])
                    if is_adjacent(c1, c2) and (vmap[c2], vmap[c1]) not in edges:
                        edges.extend([
                            (vmap[c1], vmap[c2])
                        ])


# write_to_TXT(list(vertices), list(edges), string)
# Writes the given vertices and edges to file
def write_to_TXT(vertices, edges, fileName):

    outfile = open(fileName + '.txt', 'w')

    outfile.write('%d\n' % len(vertices))
    for vertex in vertices:
        outfile.write("%d %d %d\n" % (vertex[0], vertex[1], vertex[2]))

    outfile.write('%d\n' % len(edges))
    for edge in edges:
        outfile.write('%d %d\n' % (edge[0], edge[1]))

    outfile.close()


def build(base_name, data):

    # Stores the vertex ID of a certain x, y coordinate.
    # (x, y): vertexID
    vmap = {}
    # Stores a list of vertexIDs for each square
    # (squareID): [vertexID...]
    v_sq_map = {}

    vertices = []
    edges = []

    square_list = build_squares(data['squares'])
    merge_squares(square_list, data['graph'])
    add_internal_squares(square_list, vertices, edges, data['height'])
    add_routers(square_list, vertices, edges, vmap, v_sq_map, data['height'])
    connect_adjacent_routers(square_list, vertices, edges, vmap, v_sq_map, data['height'], data['graph'])
    write_to_TXT(vertices, edges, base_name)


