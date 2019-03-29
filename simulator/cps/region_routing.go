package cps

import (
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
