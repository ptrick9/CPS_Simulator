package cps

import "math/rand"

type GridSpot struct {
	X int
	Y int
	Value int
}

type byValue []GridSpot

func (b byValue) Len() int {
	return len(b)
}

func (b byValue) Swap(i,j int) {
	b[i],b[j] = b[j],b[i]
}

func (b byValue) Less(i,j int) bool {
	return b[i].Value < b[j].Value
}

type byRandom []GridSpot

func (r byRandom) Len() int {
	return len(r)
}

func (r byRandom) Swap(i,j int) {
	r[i],r[j] = r[j],r[i]
}

func (r byRandom) Less(i,j int) bool {
	choice := rand.Intn(2)
	if choice == 0 {
		return true
	} else if choice == 1{
		return false
	}
	return true
}