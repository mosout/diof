package pso

import "math/rand"

type ParticleStatus struct {
	Fitness   float64
	Positions []float64
}

func (p *ParticleStatus) OverwriteIfBetter(ps *ParticleStatus, better func(float64, float64) bool) {
	if better(ps.Fitness, p.Fitness) {
		p.Overwrite(ps)
	}
}
func (p *ParticleStatus) Overwrite(ps *ParticleStatus) {
	p.Fitness = ps.Fitness
	if len(p.Positions) != len(ps.Positions) {
		p.Positions = make([]float64, len(ps.Positions))
	}
	copy(p.Positions, ps.Positions)
}
func (p *ParticleStatus) UpdateFitness(f func([]float64) float64) {
	p.Fitness = f(p.Positions)
}

func NewParticleStatus(dims int) ParticleStatus {
	return ParticleStatus{
		Positions: make([]float64, dims),
	}
}
func NewParticleStatusWithAnother(ps *ParticleStatus) ParticleStatus {
	p := ParticleStatus{
		Positions: make([]float64, len(ps.Positions)),
	}
	p.Overwrite(ps)
	return p
}

type Particle struct {
	Status     ParticleStatus
	BestStatus ParticleStatus
	V          []float64
	Swarm      *baseSwarm
}

func NewParticle(lowBounds []float64, intervals []float64, dims int, swarm *baseSwarm) *Particle {
	p := Particle{
		Status:     NewParticleStatus(dims),
		BestStatus: NewParticleStatus(dims),
		V:          make([]float64, dims),
		Swarm:      swarm,
	}
	for i := 0; i < dims; i++ {
		calRand := func() float64 {
			return rand.Float64() * intervals[i]
		}
		p.Status.Positions[i] = calRand() + lowBounds[i]
		p.V[i] = calRand() * 0.01
		copy(p.BestStatus.Positions, p.Status.Positions)
		p.Status.Fitness = p.Swarm.problem.F(p.Status.Positions)
		p.BestStatus.Fitness = p.Status.Fitness
	}
	return &p
}
func (p *Particle) UpdateFitness() {
	p.Status.UpdateFitness(p.Swarm.problem.F)
	p.BestStatus.OverwriteIfBetter(&p.Status, p.Swarm.better)
}
