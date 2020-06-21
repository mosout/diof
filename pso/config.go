package pso

type swarmConfig struct {
	W          float64
	C1         float64
	C2         float64
	NParticles int
}
type ServerParams struct {
	Target    string
	Namespace string
}
