package cps

import (
    "math"
)

type RoutingSquare struct {
	x1, x2, y1, y2 int
	can_cut        bool
	id_num         int
	routers        []Tuple
}

func square_equals(s1, s2 RoutingSquare) bool {
	return (s1.x1 == s2.x1 && s1.x2 == s2.x2 && s1.y1 == s2.y1 && s1.y2 == s2.y2)
}

func square_list_remove(s RoutingSquare) {
	index := -1
	for i, sq := range square_list {
		if square_equals(sq, s) {
			index = i
		}
	}
	if index != -1 {
		square_list = square_list[:index+copy(square_list[index:], square_list[index+1:])]
	}
}

func within(s RoutingSquare, p Tuple) bool {
	return (p.x >= s.x1 && p.x <= s.x2 && p.y >= s.y1 && p.y <= s.y2)
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func side_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.can_cut {
		ret := 1.1
		if sq1.x1 == sq2.x1 || sq1.x2 == sq2.x2 {
			ret = float64(sq1.x2-sq1.x1) / float64(max((sq2.x2-sq2.x1), 1))
		} else if sq1.y1 == sq2.y1 || sq1.y2 == sq2.y2 {
			ret = float64(sq1.y2-sq1.y1) / float64(max((sq2.y2-sq2.y1), 1))
		}
		if ret > 1 {
			return 0.0
		} else {
			return ret
		}
	} else {
		return 0.0
	}
}

func area_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.can_cut {
		if sq1.x1 > sq2.x1 && sq1.x2 < sq2.x2 {
			return float64(sq1.x2-sq1.x1) / float64(max((sq2.x2-sq2.x1), 1))
		} else if sq1.y1 > sq2.y1 && sq1.y2 < sq2.y2 {
			return float64(sq1.y2-sq1.y1) / float64(max((sq2.y2-sq2.y1), 1))
		} else {
			return 0.0
		}
	} else {
		return 0.0
	}
}

func single_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.x1 == sq2.x1 && sq1.x2 == sq2.x2 {
		if sq1.y1 < sq2.y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq1.y1, sq2.y2, true, sq1.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq2.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
		}
	} else if sq1.y1 == sq2.y1 && sq1.y2 == sq2.y2 {
		if sq1.x1 < sq2.x1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
		}
	} else if sq1.x1 == sq2.x1 {
		if sq1.y1 < sq2.y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq1.y1, sq2.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq2.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		}
	} else if sq1.x2 == sq2.x2 {
		if sq1.y1 < sq2.y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq1.y1, sq2.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq2.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		}
	} else if sq1.y1 == sq2.y1 {
		if sq1.x1 < sq2.x1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
		}
	} else if sq1.y2 == sq2.y2 {
		if sq1.x1 < sq2.x1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, false, sq2.id_num, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, false, sq2.id_num, make([]Tuple, 0)})
		}
	}

	return new_squares
}

func double_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.x1 > sq2.x1 && sq1.x2 < sq2.x2 {
		if sq1.y2+1 == sq2.y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq1.y1, sq2.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, false, -1, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq1.x2, sq2.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x1 - 1, sq2.y1, sq2.y2, false, sq2.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.x2 + 1, sq2.x2, sq2.y1, sq2.y2, false, -1, make([]Tuple, 0)})
		}
	} else if sq1.y1 > sq2.y1 && sq1.y2 < sq2.y2 {
		if sq1.x2+1 == sq2.x1 {
			new_squares = append(new_squares, RoutingSquare{sq1.x1, sq2.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, false, sq2.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, false, -1, make([]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq1.x2, sq1.y1, sq1.y2, true, sq1.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq2.y1, sq1.y1 - 1, false, sq2.id_num, make([]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.x1, sq2.x2, sq1.y2 + 1, sq2.y2, false, -1, make([]Tuple, 0)})
		}
	}

	return new_squares
}

func rebuild(sq_list []RoutingSquare) {
	for x := 0; x < len(square_list); x++ {
		border_dict[x] = make([]int, 0)
		square_list[x].routers = make([]Tuple, len(square_list))
	}

	for y := 0; y < len(sq_list); y++ {
		square := sq_list[y]

		for z := y + 1; z < len(sq_list); z++ {
			new_square := sq_list[z]

			if new_square.x1 >= square.x1 && new_square.x2 <= square.x2 {
				if new_square.y1 == square.y2+1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{int((new_square.x2-new_square.x1)/2) + new_square.x1, new_square.y1}
					square_list[y].routers[z] = Tuple{int((new_square.x2-new_square.x1)/2) + new_square.x1, square.y2}

				} else if new_square.y2 == square.y1-1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{int((new_square.x2-new_square.x1)/2) + new_square.x1, new_square.y2}
					square_list[y].routers[z] = Tuple{int((new_square.x2-new_square.x1)/2) + new_square.x1, square.y1}
				}
			} else if new_square.y1 >= square.y1 && new_square.y2 <= square.y2 {
				if new_square.x1 == square.x2+1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{new_square.x1, int((new_square.y2-new_square.y1)/2) + new_square.y1}
					square_list[y].routers[z] = Tuple{square.x2, int((new_square.y2-new_square.y1)/2) + new_square.y1}

				} else if new_square.x2 == square.x1-1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{new_square.x2, int((new_square.y2-new_square.y1)/2) + new_square.y1}
					square_list[y].routers[z] = Tuple{square.x1, int((new_square.y2-new_square.y1)/2) + new_square.y1}
				}
			}
			if square.x1 >= new_square.x1 && square.x2 <= new_square.x2 {
				if square.y1 == new_square.y2+1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{int((square.x2-square.x1)/2) + square.x1, square.y1}
					square_list[y].routers[z] = Tuple{int((square.x2-square.x1)/2) + square.x1, new_square.y2}

				} else if square.y2 == new_square.y1-1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{int((square.x2-square.x1)/2) + square.x1, square.y2}
					square_list[y].routers[z] = Tuple{int((square.x2-square.x1)/2) + square.x1, new_square.y1}
				}
			} else if square.y1 >= new_square.y1 && square.y2 <= new_square.y2 {
				if square.x1 == new_square.x2+1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{square.x1, int((square.y2-square.y1)/2) + square.y1}
					square_list[y].routers[z] = Tuple{new_square.x2, int((square.y2-square.y1)/2) + square.y1}

				} else if square.x2 == new_square.x1-1 {
					border_dict[y] = append(border_dict[y], z)
					border_dict[z] = append(border_dict[z], y)

					square_list[z].routers[y] = Tuple{square.x2, int((square.y2-square.y1)/2) + square.y1}
					square_list[y].routers[z] = Tuple{new_square.x1, int((square.y2-square.y1)/2) + square.y1}
				}
			}
		}
	}
}

func dist(p1, p2 Tuple) float64 {
	x_dist := math.Abs(float64(p2.x - p1.x))
	y_dist := math.Abs(float64(p2.y - p1.y))
	return math.Sqrt(float64(math.Pow(x_dist, 2) + math.Pow(y_dist, 2)))
}