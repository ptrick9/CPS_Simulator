package cps

import (
	"fmt"
	"math"
)

type Tuple struct {
	X, Y int
}

type RegionParams struct {
	Point_list []Tuple

	Point_list2 [][]bool

	Point_dict map[Tuple]bool

	Square_list []RoutingSquare

	Border_dict map[int][]int

	Node_tables []map[Tuple]float64

	Possible_paths [][]int

	Stim_list map[int]Tuple
}

func Point_list_remove(point Tuple, r *RegionParams) {
	index := -1
	for i, p := range r.Point_list {
		if (p.X == point.X) && (p.Y == point.Y) {
			index = i
		}
	}
	if index != -1 {
		r.Point_list = r.Point_list[:index+copy(r.Point_list[index:], r.Point_list[index+1:])]
	}
}

func RemoveRoutingSquare(sq RoutingSquare, r *RegionParams) {
	for i := sq.X1; i <= sq.X2; i++ {
		for j := sq.Y1; j <= sq.Y2; j++ {
			r.Point_dict[Tuple{i, j}] = false
			//r.Point_list_remove(Tuple{i, j})
			r.Point_list2[i][j] = false
		}
	}
}

func RegionContaining(p Tuple, r *RegionParams) int {
	for i, s := range r.Square_list {
		if p.X >= s.X1 && p.X <= s.X2 && p.Y >= s.Y1 && p.Y <= s.Y2 {
			return i
		}
	}
	return -1
}

func Is_in(i int, list []int) bool {
	for _, b := range list {
		if b == i {
			return true
		}
	}
	return false
}

func Search(prev_region, curr_region, end_region int, curr_path []int, r *RegionParams) {
	if curr_region == end_region {
		curr_path = append(curr_path, curr_region)
		r.Possible_paths = append(r.Possible_paths, curr_path)
	} else {
		for _, reg := range r.Border_dict[curr_region] {
			if reg == end_region {
				curr_path = append(curr_path, curr_region)
				curr_path = append(curr_path, reg)
				r.Possible_paths = append(r.Possible_paths, curr_path)
			} else if !Is_in(reg, curr_path) {
				next_path := append(curr_path, curr_region)
				Search(curr_region, reg, end_region, next_path, r)
			}
		}
	}
}

func PossPaths(p1, p2 Tuple, r *RegionParams) {
	start_region := RegionContaining(p1, r)
	end_region := RegionContaining(p2, r)

	r.Possible_paths = make([][]int, 0)

	Search(-1, start_region, end_region, make([]int, 0), r)
}

func InRegionRouting(p1, p2 Tuple) []Coord {
	ret_path := make([]Coord, 0)
	end_x := -1
	if p1.X < p2.X {
		for val := p1.X; val <= p2.X; val++ {
			ret_path = append(ret_path, Coord{X: val, Y: p1.Y})
		}
		end_x = p2.X
	} else {
		for val := p1.X; val >= p2.X; val-- {
			ret_path = append(ret_path, Coord{X: val, Y: p1.Y})
		}
		end_x = p2.X
	}
	if p1.Y < p2.Y {
		for val := p1.Y; val <= p2.Y; val++ {
			ret_path = append(ret_path, Coord{X: end_x, Y: val})
		}
	} else {
		for val := p1.Y; val >= p2.Y; val-- {
			ret_path = append(ret_path, Coord{X: end_x, Y: val})
		}
	}
	return ret_path
}

func GetPath(c1, c2 Coord, r *RegionParams) []Coord {
	p1 := Tuple{c1.X, c1.Y}
	p2 := Tuple{c2.X, c2.Y}

	PossPaths(p1, p2, r)

	min_dist := math.Pow(100, 100)
	index := -1

	for i, path := range r.Possible_paths {
		curr_dist := 0.0
		for j, region := range path {
			if len(path) == 1 {
				curr_dist += Dist(p1, p2)
			} else {
				if j == 0 {
					curr_dist += Dist(p1, r.Square_list[region].Routers[path[j+1]])
				} else if j == len(path)-1 {
					curr_dist += Dist(p2, r.Square_list[region].Routers[path[j-1]])
				} else {
					curr_dist += r.Node_tables[region][Tuple{path[j-1], path[j+1]}]
				}
			}
		}

		if curr_dist < min_dist {
			min_dist = curr_dist
			index = i
		}
	}
	ret_path := make([]Coord, 0)

	if len(r.Possible_paths[index]) == 1 {
		ret_path = append(ret_path, InRegionRouting(p1, p2)...)
	} else {
		for i, s := range r.Possible_paths[index] {
			if i == 0 {
				ret_path = append(ret_path, InRegionRouting(p1, r.Square_list[s].Routers[r.Possible_paths[index][i+1]])...)
			} else if i == len(r.Possible_paths[index])-1 {
				ret_path = append(ret_path, InRegionRouting(r.Square_list[s].Routers[r.Possible_paths[index][i-1]], p2)...)
			} else {
				ret_path = append(ret_path, InRegionRouting(r.Square_list[s].Routers[r.Possible_paths[index][i-1]], r.Square_list[s].Routers[r.Possible_paths[index][i+1]])...)
			}
		}
	}

	return ret_path
}

func GenerateRouting(p *Params, r *RegionParams) {
	fmt.Println("Beginning Region Routing")


	id_counter := 0
	done := false
	//for len(r.Point_list) != 0 {
	for !done {
		top_left := Tuple{-1, -1}
		for y := 0; y < p.Height; y++ {
			for x := 0; x < p.Width; x++ {
				//fmt.Printf("X: %d, Y: %d, v: %v progress: %d\n", x, y, r.Point_list2[x][y], len(r.Square_list))
				if r.Point_list2[x][y] {
					top_left = Tuple{x, y}
					break
				}
			}
			if (top_left != Tuple{-1, -1}) {
				break
			}
		}
		//fmt.Printf("working %d %d\n", top_left.X, top_left.Y)
		if (top_left == Tuple{-1, -1}) {
			done = true
			break
		}
		//top_left := r.Point_list[0]
		temp := Tuple{top_left.X, top_left.Y}

		for r.Point_dict[Tuple{temp.X + 1, temp.Y}] {
			temp.X += 1
			//fmt.Printf("dict: X: %d, Y: %d", temp.X+1, temp.Y)
		}

		collide := false
		y_test := Tuple{top_left.X, top_left.Y}

		for !collide {
			y_test.Y += 1

			for x_val := top_left.X; x_val < temp.X; x_val++ {
				if !r.Point_dict[Tuple{x_val, y_test.Y}] {
					collide = true
				}
			}
		}

		bottom_right := Tuple{temp.X, y_test.Y - 1}

		//intln(top_left.X, bottom_right.X, top_left.Y, bottom_right.Y)

		new_square := RoutingSquare{top_left.X, bottom_right.X, top_left.Y, bottom_right.Y, true, id_counter, make([]Tuple, 0)}
		id_counter++
		//fmt.Println("start_r_square")
		RemoveRoutingSquare(new_square, r)
		//fmt.Println("end_r_square")
		r.Square_list = append(r.Square_list, new_square)
	}
	//fmt.Println("Built all squares")
	length := len(r.Square_list)
	for y, _ := range r.Square_list {
		square := r.Square_list[y]
		r.Square_list[y].Routers = make([]Tuple, length)

		for z := y + 1; z < len(r.Square_list); z++ {
			new_square := r.Square_list[z]

			if new_square.X1 >= square.X1 && new_square.X2 <= square.X2 {
				if new_square.Y1 == square.Y2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

				} else if new_square.Y2 == square.Y1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)
				}
			} else if new_square.Y1 >= square.Y1 && new_square.Y2 <= square.Y2 {
				if new_square.X1 == square.X2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

				} else if new_square.X2 == square.X1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)
				}
			}
			if square.X1 >= new_square.X1 && square.X2 <= new_square.X2 {
				if square.Y1 == new_square.Y2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

				} else if square.Y2 == new_square.Y1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)
				}
			} else if square.Y1 >= new_square.Y1 && square.Y2 <= new_square.Y2 {
				if square.X1 == new_square.X2+1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)

				} else if square.X2 == new_square.X1-1 {
					r.Border_dict[y] = append(r.Border_dict[y], z)
					r.Border_dict[z] = append(r.Border_dict[z], y)
				}
			}
		}
	}
	//fmt.Println(r.Border_dict)

	//Cutting takes place in this loop
	for true {
		rebuilt := false

		for i := 0; i < len(r.Square_list) && !rebuilt; i++ {

			for _, n := range r.Border_dict[i] {

				s_rat := Side_ratio(r.Square_list[i], r.Square_list[n])
				if s_rat > 0.6 {
					new_squares := Single_cut(r.Square_list[i], r.Square_list[n])

					s1 := r.Square_list[n]
					s2 := r.Square_list[i]

					Square_list_remove(s1, r)
					Square_list_remove(s2, r)

					r.Square_list = append(r.Square_list, new_squares...)

					Rebuild(r.Square_list, r)

					rebuilt = true

					break
				}

				a_rat := Area_ratio(r.Square_list[i], r.Square_list[n])
				if a_rat > 0.6 {
					new_squares := Double_cut(r.Square_list[i], r.Square_list[n])

					s1 := r.Square_list[n]
					s2 := r.Square_list[i]

					Square_list_remove(s1, r)
					Square_list_remove(s2, r)

					new_squares[2].Id_num = len(r.Square_list)

					r.Square_list = append(r.Square_list, new_squares...)

					Rebuild(r.Square_list, r)

					rebuilt = true

					break
				}
			}
		}

		if !rebuilt {
			break
		}
	}

	r.Node_tables = make([]map[Tuple]float64, len(r.Square_list))

	for key, values := range r.Border_dict {
		if key < len(r.Square_list) {
			r.Node_tables[key] = make(map[Tuple]float64)
			if len(values) > 1 {
				for n := 0; n < len(values); n++ {
					next := n + 1
					for next < len(values) {
						node_a := r.Border_dict[key][n]
						node_b := r.Border_dict[key][next]

						p1 := r.Square_list[key].Routers[node_a]
						p2 := r.Square_list[key].Routers[node_b]

						r.Node_tables[key][Tuple{node_a, node_b}] = Dist(p1, p2)
						r.Node_tables[key][Tuple{node_b, node_a}] = Dist(p1, p2)

						next += 1
					}
				}
			}
		}
	}

}
