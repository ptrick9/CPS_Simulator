package main

import (
	"math"
)

func point_list_remove(point Tuple) {
	index := -1
	for i, p := range point_list {
		if (p.x == point.x) && (p.y == point.y) {
			index = i
		}
	}
	if index != -1 {
		point_list = point_list[:index+copy(point_list[index:], point_list[index+1:])]
	}
}

func removeRoutingSquare(sq RoutingSquare) {
	for i := sq.x1; i <= sq.x2; i++ {
		for j := sq.y1; j <= sq.y2; j++ {
			point_dict[Tuple{i, j}] = false
			//point_list_remove(Tuple{i, j})
			point_list2[i][j] = false
		}
	}
}

func regionContaining(p Tuple) int {
	for i, s := range square_list {
		if p.x >= s.x1 && p.x <= s.x2 && p.y >= s.y1 && p.y <= s.y2 {
			return i
		}
	}
	return -1
}

func is_in(i int, list []int) bool {
	for _, b := range list {
		if b == i {
			return true
		}
	}
	return false
}

func search(prev_region, curr_region, end_region int, curr_path []int) {
	if curr_region == end_region {
		curr_path = append(curr_path, curr_region)
		possible_paths = append(possible_paths, curr_path)
	} else {
		for _, r := range border_dict[curr_region] {
			if r == end_region {
				curr_path = append(curr_path, curr_region)
				curr_path = append(curr_path, r)
				possible_paths = append(possible_paths, curr_path)
			} else if !is_in(r, curr_path) {
				next_path := append(curr_path, curr_region)
				search(curr_region, r, end_region, next_path)
			}
		}
	}
}

func possPaths(p1, p2 Tuple) {
	start_region := regionContaining(p1)
	end_region := regionContaining(p2)

	possible_paths = make([][]int, 0)

	search(-1, start_region, end_region, make([]int, 0))
}

func inRegionRouting(p1, p2 Tuple) []Coord {
	ret_path := make([]Coord, 0)
	end_x := -1
	if p1.x < p2.x {
		for val := p1.x; val <= p2.x; val++ {
			ret_path = append(ret_path, Coord{x: val, y: p1.y})
		}
		end_x = p2.x
	} else {
		for val := p1.x; val >= p2.x; val-- {
			ret_path = append(ret_path, Coord{x: val, y: p1.y})
		}
		end_x = p2.x
	}
	if p1.y < p2.y {
		for val := p1.y; val <= p2.y; val++ {
			ret_path = append(ret_path, Coord{x: end_x, y: val})
		}
	} else {
		for val := p1.y; val >= p2.y; val-- {
			ret_path = append(ret_path, Coord{x: end_x, y: val})
		}
	}
	return ret_path
}

func getPath(c1, c2 Coord) []Coord {
	p1 := Tuple{c1.x, c1.y}
	p2 := Tuple{c2.x, c2.y}

	possPaths(p1, p2)

	min_dist := math.Pow(100, 100)
	index := -1

	for i, path := range possible_paths {
		curr_dist := 0.0
		for j, region := range path {
			if len(path) == 1 {
				curr_dist += dist(p1, p2)
			} else {
				if j == 0 {
					curr_dist += dist(p1, square_list[region].routers[path[j+1]])
				} else if j == len(path)-1 {
					curr_dist += dist(p2, square_list[region].routers[path[j-1]])
				} else {
					curr_dist += node_tables[region][Tuple{path[j-1], path[j+1]}]
				}
			}
		}

		if curr_dist < min_dist {
			min_dist = curr_dist
			index = i
		}
	}
	ret_path := make([]Coord, 0)

	if len(possible_paths[index]) == 1 {
		ret_path = append(ret_path, inRegionRouting(p1, p2)...)
	} else {
		for i, s := range possible_paths[index] {
			if i == 0 {
				ret_path = append(ret_path, inRegionRouting(p1, square_list[s].routers[possible_paths[index][i+1]])...)
			} else if i == len(possible_paths[index])-1 {
				ret_path = append(ret_path, inRegionRouting(square_list[s].routers[possible_paths[index][i-1]], p2)...)
			} else {
				ret_path = append(ret_path, inRegionRouting(square_list[s].routers[possible_paths[index][i-1]], square_list[s].routers[possible_paths[index][i+1]])...)
			}
		}
	}

	return ret_path
}
