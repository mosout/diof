package pso

func SelectBetter(target string) func(float64, float64) bool {
	if target == "min" {
		return func(x, y float64) bool {
			return x < y
		}
	}
	if target == "max" {
		return func(x, y float64) bool {
			return x > y
		}
	}
	return nil
}
