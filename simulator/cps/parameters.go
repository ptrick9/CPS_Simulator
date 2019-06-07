package cps

import (
	"math/rand"
	"time"
)

/*This Go file deals with properly set parameters of the nodes
*/

var(
	full = 318 //Custom distribution array size
	vardym = [][]int{ // Custom distribution array
		{87,18},{88,17},{86,17},{89,16},{85,16},{90,15},{84,15},{91,14},
		{83,14},{92,13},{82,13},{93,12},{81,12},{94,11},{80,11},{95,10},
		{79,10},{96,9},{78,9},{97,8},{77,8},{98,7},{76,7},{99,6},{75,6},
		{100,5},{74,5},{73,4},{72,3},{71,2},{70,1},{60,1},{50,1},{30,1},
		{20,1}}
)

//The 2 functions are what we will most likely use
//produces a normal distribution sample with specific sigma and mu
func GetNormDistro(numNodes int, mu float32, sigma float32) (arr []float32) {
	x := float32(0)
	for true {
		x = NormalInverse(mu,sigma)
		if x <= 100 && x >= 0 {
			arr = append(arr,x)
			if len(arr) >= numNodes {
				return
			}
		}
	}
	return
}
//produces normal sample
func NormalInverse(mu float32, sigma float32) (float32) {
	return float32(rand.NormFloat64() * float64(sigma) + float64(mu))
}


//Returns linear list of battery values, very good for debugging because it is in order
func GetLinearBatteryValues(numNodes int) (y []float32) {
	step := 100.0 / (float32(numNodes))
	for numNodes >= 1 {
		y = append(y, step * (float32(numNodes)))
		numNodes -= 1
	}
	return
}

//Returns a constant value for a uniform battery loss for all the nodes
func GetLinearBatteryLossConstant(numNodes int, lossConst float32) (y []float32) {
	for numNodes >= 1 {
		y = append(y,lossConst)
		numNodes -= 1
	}
	return
}
//Created custom distribution
func ProduceCustomDistribution(numb int) (y []float32) {
	//Traces determine how many of each variable one gets from their range
	traceRise := 18
	traceFall := 18
	//Ranges are the variable ranges one puts in their distribution
	//range 1 ramps up
	rangeStart1 := 87
	rangeFinish1 := 100
	//range 2 ramps down
	rangeStart2 := 86
	rangeFinish2 := 70
	//rising distribution
	for i := rangeStart1; i <= rangeFinish1; i++ {
		for j := traceRise; j >= 1; j-- {
			y = append(y,float32(i))
		}
	}
	for i := rangeStart2; i >= rangeFinish2; i-- {
		for j := 1; j >= traceFall; j-- {
			y = append(y, float32(i))
		}
	}
	//add odds and ends here
	return
}

//returns charges from custom distribution array
func GetinitialChargeDynamic(numNodes int) (y []float32) {
	for true {
		for i := 0; i < len(vardym); i++ {
			for j := 0; j < vardym[i][1]; j++ {
				y = append(y, float32(vardym[i][0]))
				numNodes -= 1
				if numNodes == 0 {
					return Shuffle(y)
				}
			}
		}
	}
	return
}

//returns charges from custom distribution array ratio to scale of original array
func GetInitialChargeSuperDynamic(numNodes int, scalar float32) (y []float32) {
	for i := 0; i < len(vardym); i++ {
		for j := 0; j < vardym[i][1]; j++ {
			y = append(y,float32(vardym[i][0]) * scalar)
		}
	}
	if numNodes < full {
		for i := 0; i < full-numNodes; i++ {
			y = RemoveFloat32(y,RandomInt(0,len(y)-1))
		}
	} else if numNodes > full {
		for i := 0; i < numNodes-full; i++ {
			y = append(y,y[RandomInt(0,len(y)-1)])
		}
	}
	return Shuffle(y)
}

//returns index of int in array
func GetIndexInt(slice []int, element int) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

//removes index from array
func RemoveFloat32(slice []float32, s int) []float32 {
	return append(slice[:s],slice[s+1:]...)
}
//returns random number between 2 numbers
func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
//shuffles array
func Shuffle(a []float32) ([]float32) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for n := len(a); n > 0; n-- {
		randIndex := r.Intn(n)
		a[n-1], a[randIndex] = a[randIndex],a[n-1]
	}
	return a
}