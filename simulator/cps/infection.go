package cps

//
func calculateAlpha(min float32, max float32, target float32) float32 {
	return (target - min) / (max - min)
}

func lerp(min float32, max float32, alpha float32) float32 {
	return min + (max - min) * alpha
}
