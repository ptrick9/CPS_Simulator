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

	Checked		[]int
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
		if p.X >= s.X1 && p.X <= s.X2 && p.Y <= s.Y1 && p.Y >= s.Y2 {
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

func ValidPath(reg int, endpoint Coord, r *RegionParams) bool{
	end := RegionContaining(Tuple{endpoint.X, endpoint.Y}, r)
	if reg == end {
		r.Checked = make([]int,0)
		return true
	}
	if len(r.Border_dict[reg]) == 0 {
		r.Checked = make([]int,0)
		return false
	} else {
		for i := 0; i < len(r.Border_dict[reg]); i++ {
			if r.Border_dict[reg][i] == end {
				//fmt.Printf("Found a path to region %v\n",r.Border_dict[reg][i])
				r.Checked = make([]int,0)
				return true
			}
			if r.Border_dict[reg][i] != reg && !Is_in(r.Border_dict[reg][i], r.Checked){
				r.Checked = append(r.Checked, reg)
				//fmt.Println(r.Checked)
				if ValidPath(r.Border_dict[reg][i], endpoint, r) {
					return true
				}
			}
		}
	}
	return false
}


func PossPaths(p1, p2 Tuple, r *RegionParams) {
	start_region := RegionContaining(p1, r)
	end_region := RegionContaining(p2, r)

	r.Possible_paths = make([][]int, 0)

	if ValidPath(start_region, Coord{X:p2.X, Y:p2.Y}, r) {
		Search(-1, start_region, end_region, make([]int, 0), r)
		//fmt.Printf("%v %v %v\n", start_region, end_region, r.Possible_paths)
	}
	//Search(-1, start_region, end_region, make([]int, 0), r)
}

func InRegionRouting(p1, p2 Tuple, r *RegionParams) []Coord {
	/*square1 := r.Square_list[RegionContaining(p1, r)]
	square2 := r.Square_list[RegionContaining(p2, r)]
	//fmt.Printf("Region containing %v is %v\n", p1, RegionContaining(p1, r))
	//fmt.Printf("Region containing %v is %v\n", p2, RegionContaining(p2, r))

	val1_first := -1
	val1_second := -1
	val2_first := -1
	val2_second := -1
	if square2.X2 < square1.X1 || square2.X1 > square1.X2 {
		val1_first = p1.Y
		val1_second = p1.X
		val2_first = p2.Y
		val2_second = p2.X
	}
	if square2.Y2 < square1.Y1 || square2.Y1 > square1.Y2 {
		val1_first = p1.X
		val1_second = p1.Y
		val2_first = p2.X
		val2_second = p2.Y
	}

	ret_path := make([]Coord, 0)
	end_x := -1
	if val1_first < val2_first {
		//fmt.Println("Moving right")
		for val := val1_first; val <= val2_first; val++ {
			//fmt.Printf("(%v,%v)\n",val,val1_second)
			ret_path = append(ret_path, Coord{X: val, Y: val1_second})
		}
		end_x = val2_first
	} else {
		//fmt.Println("Moving left")
		for val := val1_first; val >= val2_first; val-- {
			//fmt.Printf("(%v,%v)\n",val,val1_second)
			ret_path = append(ret_path, Coord{X: val, Y: val1_second})
		}
		end_x = val2_first
	}
	if val1_second < val2_second {
		//fmt.Println("Moving up")
		for val := val1_second; val <= val2_second; val++ {
			//fmt.Printf("(%v,%v)\n",end_x,val)
			ret_path = append(ret_path, Coord{X: end_x, Y: val})
		}
	} else {
		//fmt.Println("Moving down")
		for val := val1_second; val >= val2_second; val-- {
			//fmt.Printf("(%v,%v)\n",end_x,val)
			ret_path = append(ret_path, Coord{X: end_x, Y: val})
		}
	}
	fmt.Printf("Path Taken: %v\n", ret_path)
	return ret_path*/

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

func GetPath(c1, c2 Coord, r *RegionParams, p *Params) []Coord {
	//fmt.Println(r.Possible_paths)
	p1 := Tuple{c1.X, c1.Y}
	p2 := Tuple{c2.X, c2.Y}
	PossPaths(p1, p2, r)

	min_dist := math.Pow(100, 100)
	index := -1
	index2 := -1
	//fmt.Printf("%v %v %v\n", c1, c2, r.Possible_paths)
	for i, path := range r.Possible_paths {
		for k := 0; k < 3; k++ {
			curr_dist := 0.0
			for j, region := range path {
					if len(path) == 1 {
						curr_dist += Dist(p1, p2)
					} else {
						if j == 0 {
							curr_dist += Dist(p1, r.Square_list[region].Routers[path[j+1]][k])
						} else if j == len(path)-1 {
							curr_dist += Dist(p2, r.Square_list[region].Routers[path[j-1]][k])
						} else {
							//curr_dist += r.Node_tables[region][Tuple{path[j-1], path[j+1]}]
							curr_dist += Dist(p1, r.Square_list[region].Routers[path[j]][k])
						}
					}
				}
				//fmt.Printf("Current distance: %v, Path: %v\n", curr_dist, k)
				if curr_dist < min_dist {
					min_dist = curr_dist
					//fmt.Printf("Current min distance: %v, Taking path: %v\n", min_dist, k)
					index = i
					index2 = k
				}
		}
	}
	ret_path := make([]Coord, 0)

	//if r.Possible_paths[index] == 5

	if len(r.Possible_paths[index]) == 1 {
		//fmt.Println("Possible paths length 1")
		ret_path = append(ret_path, InRegionRouting(p1, p2, r)...)
	} else {
		for i, s := range r.Possible_paths[index] {
				if i == 0 {
					ret_path = append(ret_path, InRegionRouting(p1, r.Square_list[s].Routers[r.Possible_paths[index][i+1]][index2],r)...)
				} else if i == len(r.Possible_paths[index])-1 {
					ret_path = append(ret_path, InRegionRouting(r.Square_list[s].Routers[r.Possible_paths[index][i-1]][index2], p2,r)...)
				} else {
					//fmt.Println(r.Square_list[s].Routers[r.Possible_paths[index][i-1]][index2])
					ret_path = append(ret_path, InRegionRouting(r.Square_list[s].Routers[r.Possible_paths[index][i-1]][index2], r.Square_list[s].Routers[r.Possible_paths[index][i+1]][index2],r)...)
				}
		}
		/*start_point := p1
		end_point := Tuple{X: 0, Y: 0}

		for i, s := range r.Possible_paths[index] {
			if i == len(r.Possible_paths[index])-1 {
				//fmt.Println("Possible paths[index] length - 1 = i")
				ret_path = append(ret_path, InRegionRouting(start_point, p2,r)...)
			} else {
				end_point = ClosestToSquare(start_point, r.Square_list[s], p)
				fmt.Printf("Start point: %v, End point after :%v\n", start_point, end_point)
				ret_path = append(ret_path, InRegionRouting(start_point, end_point, r)...)
				start_point = end_point
			}
		}*/
	}
	fmt.Println(ret_path)
	var tmp Coord = ret_path[0]
	for i := range ret_path {
		if ret_path[i].X == -1 || ret_path[i].Y == -1 {
			fmt.Println("Error in path finding!")
		}
		if ret_path[i].X - 1 > tmp.X || ret_path[i].X + 1 < tmp.X {
			fmt.Println("Error in path finding! Node teleported on X")
		}
		if ret_path[i].Y - 1 > tmp.Y || ret_path[i].Y + 1 < tmp.Y {
			fmt.Println("Error in path finding! Node teleported on Y")
		}
		tmp = ret_path[i]
	}
	return ret_path
}

func ClosestToSquare(t Tuple, r RoutingSquare, p *Params) Tuple {
	dist := float64(p.MaxX * p.MaxY)
	newX := -1
	newY := -1

	//Loops through the list to find the closest Coord
	for x := r.X1; x < r.X2; x++ {
		for y := r.Y2; y < r.Y1; y++ {
			p := Coord{X: x, Y: y}
			newDist := math.Sqrt(math.Pow(float64(p.X-t.X), 2.0) + math.Pow(float64(p.Y-t.Y), 2.0))

			//Saves the value of that smallest distance to return
			if newDist < dist {
				dist = newDist
				newX = x
				newY = y
			}
		}
	}
	return Tuple{X: newX, Y: newY}
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

		//fmt.Println(top_left.X, bottom_right.X, top_left.Y, bottom_right.Y)

		new_square := RoutingSquare{top_left.X, bottom_right.X, top_left.Y, bottom_right.Y, true, id_counter, make([][]Tuple, 0)}
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
		r.Square_list[y].Routers = make([][]Tuple, length)

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
		for i := 0; i < 3; i++ {
			if key < len(r.Square_list) {
				r.Node_tables[key] = make(map[Tuple]float64)
				if len(values) > 1 {
					for n := 0; n < len(values); n++ {
						next := n + 1
						for next < len(values) {
							node_a := r.Border_dict[key][n]
							node_b := r.Border_dict[key][next]

							p1 := r.Square_list[key].Routers[node_a][i]
							p2 := r.Square_list[key].Routers[node_b][i]

							r.Node_tables[key][Tuple{node_a, node_b}] = Dist(p1, p2)
							r.Node_tables[key][Tuple{node_b, node_a}] = Dist(p1, p2)

							next += 1
						}
					}
				}
			}
		}
	}

}
