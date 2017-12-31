package suggestor

import "math"

func CalcSampleSize(total int) float64 {

	var n float64

	// Standard Deviation
	var o float64 = 0.5

	// Level of trustworthiness 99% high value 2.58 and 95% min value 1.96
	var z float64 = 2.58

	// Limit error acceptable from 1% to 9% . 5% is value standard
	var e float64 = 0.5

	var N float64 = float64(total)

	// Formula to calc representative sample
	n = (math.Pow(z, 2)*math.Pow(o, 2)*N)/(math.Pow(e, 2)*(N-1)+math.Pow(z, 2)*math.Pow(o, 2))

	return n*100
}
