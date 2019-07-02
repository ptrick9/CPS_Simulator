package cps

import (
	"math"
)

//Square generated for pathfinding
type RoutingSquare struct {
	X1, X2, Y1, Y2 int
	Can_cut        bool
	Id_num         int
	Routers        [][]Tuple
}

//Square_equals is  comparator for two RoutingSquares
func Square_equals(s1, s2 RoutingSquare) bool {
	return (s1.X1 == s2.X1 && s1.X2 == s2.X2 && s1.Y1 == s2.Y1 && s1.Y2 == s2.Y2)
}

//Square_list_remove finds duplicate squares and removes them from the square list
func Square_list_remove(s RoutingSquare, r *RegionParams) {
	index := -1
	for i, sq := range r.Square_list {
		if Square_equals(sq, s) {
			index = i
		}
	}
	if index != -1 {
		r.Square_list = r.Square_list[:index+copy(r.Square_list[index:], r.Square_list[index+1:])]
	}
}

//Within returns whether or not a set of coordinates is within a square
func Within(s RoutingSquare, p Tuple) bool {
	return (p.X >= s.X1 && p.X <= s.X2 && p.Y >= s.Y1 && p.Y <= s.Y2)
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func Side_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.Can_cut {
		ret := 1.1
		if sq1.X1 == sq2.X1 || sq1.X2 == sq2.X2 {
			ret = float64(sq1.X2-sq1.X1) / float64(Max((sq2.X2-sq2.X1), 1))
		} else if sq1.Y1 == sq2.Y1 || sq1.Y2 == sq2.Y2 {
			ret = float64(sq1.Y2-sq1.Y1) / float64(Max((sq2.Y2-sq2.Y1), 1))
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

func Area_ratio(sq1, sq2 RoutingSquare) float64 {
	if sq1.Can_cut {
		if sq1.X1 > sq2.X1 && sq1.X2 < sq2.X2 {
			return float64(sq1.X2-sq1.X1) / float64(Max((sq2.X2-sq2.X1), 1))
		} else if sq1.Y1 > sq2.Y1 && sq1.Y2 < sq2.Y2 {
			return float64(sq1.Y2-sq1.Y1) / float64(Max((sq2.Y2-sq2.Y1), 1))
		} else {
			return 0.0
		}
	} else {
		return 0.0
	}
}

func Single_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.X1 == sq2.X1 && sq1.X2 == sq2.X2 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
		}
	} else if sq1.Y1 == sq2.Y1 && sq1.Y2 == sq2.Y2 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
		}
	} else if sq1.X1 == sq2.X1 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		}
	} else if sq1.X2 == sq2.X2 {
		if sq1.Y1 < sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		}
	} else if sq1.Y1 == sq2.Y1 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
		}
	} else if sq1.Y2 == sq2.Y2 {
		if sq1.X1 < sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([][]Tuple, 0)})
		}
	}

	return new_squares
}

func Double_cut(sq1, sq2 RoutingSquare) []RoutingSquare {
	new_squares := make([]RoutingSquare, 0)

	if sq1.X1 > sq2.X1 && sq1.X2 < sq2.X2 {
		if sq1.Y2+1 == sq2.Y1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq1.Y1, sq2.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, -1, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq1.X2, sq2.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X1 - 1, sq2.Y1, sq2.Y2, false, sq2.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq1.X2 + 1, sq2.X2, sq2.Y1, sq2.Y2, false, -1, make([][]Tuple, 0)})
		}
	} else if sq1.Y1 > sq2.Y1 && sq1.Y2 < sq2.Y2 {
		if sq1.X2+1 == sq2.X1 {
			new_squares = append(new_squares, RoutingSquare{sq1.X1, sq2.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, -1, make([][]Tuple, 0)})
		} else {
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq1.X2, sq1.Y1, sq1.Y2, true, sq1.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq2.Y1, sq1.Y1 - 1, false, sq2.Id_num, make([][]Tuple, 0)})
			new_squares = append(new_squares, RoutingSquare{sq2.X1, sq2.X2, sq1.Y2 + 1, sq2.Y2, false, -1, make([][]Tuple, 0)})
		}
	}

	return new_squares
}

func Rebuild(sq_list []RoutingSquare, r *RegionParams) {
	for x := 0; x < len(r.Square_list); x++ {
		r.Border_dict[x] = make([]int, 0)
		r.Square_list[x].Routers = make([][]Tuple, len(r.Square_list))
		for i:= range r.Square_list[x].Routers {
			r.Square_list[x].Routers[i] = make([]Tuple, 3)
		}
	}

	for y := 0; y < len(sq_list); y++ {
		square := sq_list[y]

		for z := y + 1; z < len(sq_list); z++ {
			new_square := sq_list[z]

			if new_square.X1 >= square.X1 && new_square.X2 <= square.X2 {
				if new_square.Y1 == square.Y2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{int(new_square.X1), new_square.Y1}
					r.Square_list[y].Routers[z][0] = Tuple{int(new_square.X1), square.Y2}

					r.Square_list[z].Routers[y][1] = Tuple{int((new_square.X2-new_square.X1)/2) + new_square.X1, new_square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{int((new_square.X2-new_square.X1)/2) + new_square.X1, square.Y2}

					r.Square_list[z].Routers[y][2] = Tuple{int(new_square.X2), new_square.Y1}
					r.Square_list[y].Routers[z][2] = Tuple{int(new_square.X2), square.Y2}

				} else if new_square.Y2 == square.Y1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{int(new_square.X1), new_square.Y2}
					r.Square_list[y].Routers[z][0] = Tuple{int(new_square.X1), square.Y1}

					r.Square_list[z].Routers[y][1] = Tuple{int((new_square.X2-new_square.X1)/2) + new_square.X1, new_square.Y2}
					r.Square_list[y].Routers[z][1] = Tuple{int((new_square.X2-new_square.X1)/2) + new_square.X1, square.Y1}

					r.Square_list[z].Routers[y][2] = Tuple{int(new_square.X2), new_square.Y2}
					r.Square_list[y].Routers[z][2] = Tuple{int(new_square.X2), square.Y1}
				}
			} else if new_square.Y1 >= square.Y1 && new_square.Y2 <= square.Y2 {
				if new_square.X1 == square.X2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{new_square.X1, int(new_square.Y1)}
					r.Square_list[y].Routers[z][0] = Tuple{square.X2, int(new_square.Y1)}

					r.Square_list[z].Routers[y][1] = Tuple{new_square.X1, int((new_square.Y2-new_square.Y1)/2) + new_square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{square.X2, int((new_square.Y2-new_square.Y1)/2) + new_square.Y1}

					r.Square_list[z].Routers[y][2] = Tuple{new_square.X1, int(new_square.Y2)}
					r.Square_list[y].Routers[z][2] = Tuple{square.X2, int(new_square.Y2)}

				} else if new_square.X2 == square.X1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{new_square.X2, int(new_square.Y1)}
					r.Square_list[y].Routers[z][0] = Tuple{square.X1, int(new_square.Y1)}

					r.Square_list[z].Routers[y][1] = Tuple{new_square.X2, int((new_square.Y2-new_square.Y1)/2) + new_square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{square.X1, int((new_square.Y2-new_square.Y1)/2) + new_square.Y1}

					r.Square_list[z].Routers[y][2] = Tuple{new_square.X2, int(new_square.Y2)}
					r.Square_list[y].Routers[z][2] = Tuple{square.X1, int(new_square.Y2)}
				}
			}
			if square.X1 >= new_square.X1 && square.X2 <= new_square.X2 {
				if square.Y1 == new_square.Y2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{int(square.X1), square.Y1}
					r.Square_list[y].Routers[z][0] = Tuple{int(square.X1), new_square.Y2}

					r.Square_list[z].Routers[y][1] = Tuple{int((square.X2-square.X1)/2) + square.X1, square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{int((square.X2-square.X1)/2) + square.X1, new_square.Y2}

					r.Square_list[z].Routers[y][2] = Tuple{int(square.X2), square.Y1}
					r.Square_list[y].Routers[z][2] = Tuple{int(square.X2), new_square.Y2}

				} else if square.Y2 == new_square.Y1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{int(square.X1), square.Y2}
					r.Square_list[y].Routers[z][0] = Tuple{int(square.X1), new_square.Y1}

					r.Square_list[z].Routers[y][1] = Tuple{int((square.X2-square.X1)/2) + square.X1, square.Y2}
					r.Square_list[y].Routers[z][1] = Tuple{int((square.X2-square.X1)/2) + square.X1, new_square.Y1}

					r.Square_list[z].Routers[y][1] = Tuple{int(square.X2), square.Y2}
					r.Square_list[y].Routers[z][1] = Tuple{int(square.X2), new_square.Y1}
				}
			} else if square.Y1 >= new_square.Y1 && square.Y2 <= new_square.Y2 {
				if square.X1 == new_square.X2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{square.X1, int(square.Y1)}
					r.Square_list[y].Routers[z][0] = Tuple{new_square.X2, int(square.Y1)}

					r.Square_list[z].Routers[y][1] = Tuple{square.X1, int((square.Y2-square.Y1)/2) + square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{new_square.X2, int((square.Y2-square.Y1)/2) + square.Y1}

					r.Square_list[z].Routers[y][2] = Tuple{square.X1, int(square.Y2)}
					r.Square_list[y].Routers[z][2] = Tuple{new_square.X2, int(square.Y2)}

				} else if square.X2 == new_square.X1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

					r.Square_list[z].Routers[y][0] = Tuple{square.X2, int(square.Y1)}
					r.Square_list[y].Routers[z][0] = Tuple{new_square.X1, int(square.Y1)}

					r.Square_list[z].Routers[y][1] = Tuple{square.X2, int((square.Y2-square.Y1)/2) + square.Y1}
					r.Square_list[y].Routers[z][1] = Tuple{new_square.X1, int((square.Y2-square.Y1)/2) + square.Y1}

					r.Square_list[z].Routers[y][2] = Tuple{square.X2, int(square.Y2)}
					r.Square_list[y].Routers[z][2] = Tuple{new_square.X1, int(square.Y2)}
				}
			}
		}
	}
}

func Dist(p1, p2 Tuple) float64 {
	x_dist := math.Abs(float64(p2.X - p1.X))
	y_dist := math.Abs(float64(p2.Y - p1.Y))
	return math.Sqrt(float64(math.Pow(x_dist, 2) + math.Pow(y_dist, 2)))
}
