package pso

type GlobalBestSwarm struct {
	baseSwarm
}

func NewGlobalBestSwarm(w float64, c1 float64, c2 float64, nParticles int) *GlobalBestSwarm {
	return &GlobalBestSwarm{
		baseSwarm: newBaseSwarm(w, c1, c2, nParticles),
	}
}

func (s *GlobalBestSwarm) run(numIterations int, f func(left chan *ParticleStatus, right chan *ParticleStatus)) chan float64 {
	fitnessChan := make(chan float64, numIterations)
	leftChan, rightChan := make(chan *ParticleStatus, 1), make(chan *ParticleStatus, 1)
	go func() {
		bestStatus := NewParticleStatusWithAnother(&s.particles[0].Status)
		for i := 0; i < numIterations; i++ {
			for j := 0; j < s.config.NParticles; j++ {
				bestStatus.OverwriteIfBetter(&s.particles[j].Status, s.better)
			}
			fitnessChan <- bestStatus.Fitness
			leftChan <- &bestStatus
			f(leftChan, rightChan)
			newBestStatus := <-rightChan
			for j := 0; j < s.config.NParticles; j++ {
				for k := 0; k < s.problem.Dims(); k++ {
					s.UpdateParticleVAndPosition(s.particles[j], k, newBestStatus.Positions[k])
				}
				s.particles[j].UpdateFitness()
			}
		}
		close(fitnessChan)
		close(rightChan)
		close(leftChan)
		s.stopCh <- true
	}()
	s.isRunning = true
	return fitnessChan
}
func (s *GlobalBestSwarm) Run(numIterations int) chan float64 {
	if s.rpcClient == nil {
		return s.run(numIterations, func(left chan *ParticleStatus, right chan *ParticleStatus) {
			right <- <-left
		})
	} else {
		return s.run(numIterations, func(left chan *ParticleStatus, right chan *ParticleStatus) {
			s.Push(<-left)
			right <- s.Pull()
		})
	}
}
