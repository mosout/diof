# Diof

Diof is a distributed intelligent optimization framework.

## Installation

To install Diof, run this command in your terminal:

```shell
$ go get -u github.com/mosout/diof
```

### Basic Usage
#### Define your problem
First, you need to define a struct to describe the problem you want to solve.
This struct should implement all the methods in the interface `problem.IProblem`.
```golang
type IProblem interface {
	Bounds() [][2]float64
	Dims() int
	F([]float64) float64
	Target() string
}
```
For example:
```golang
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
```
If you want to find the minimum value,`Target()` should return `min`.
Similarly if you want to find the maximum value,`Target()` should return `max. `
#### Solve your problem
Now,let us solve your problem with `GlobalBestSwarm`:
```golang
func main() {
	rand.Seed(time.Now().UnixNano())
	problem := problems.NewRAS(10, 5.0)
	s := pso.NewGlobalBestSwarm(0.5, 2, 2, 100)
	s.BindProblem(problem)
	fitnessChan := s.Run(500)
	for fitness := range fitnessChan {
		fmt.Println(fitness)
	}
	defer s.Wait()
}
```