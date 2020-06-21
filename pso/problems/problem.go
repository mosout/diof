package problems

type IProblem interface {
	Bounds() [][2]float64
	Dims() int
	F([]float64) float64
	Target() string
}
