package problems

import "math"

type RAS struct {
	dims   int
	k      float64
	bounds [][2]float64
}

func NewRAS(dims int, k float64) *RAS {
	return &RAS{
		dims: dims,
		k:    k,
	}
}
func (r *RAS) Dims() int {
	return r.dims
}
func (r *RAS) Target() string {
	return "min"
}
func (r *RAS) F(x []float64) float64 {
	result := 10.0 * float64(len(x))
	for _, v := range x {
		result += math.Pow(v, 2) - 10*math.Cos(2*math.Pi*v)
	}
	return result
}
func (r *RAS) Bounds() [][2]float64 {
	if r.bounds == nil {
		r.bounds = make([][2]float64, r.dims)
		for i := 0; i < r.dims; i++ {
			r.bounds[i] = [2]float64{-r.k, r.k}
		}
	}
	return r.bounds
}
