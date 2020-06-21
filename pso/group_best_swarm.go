package pso

type GroupBestSwarm struct {
	baseSwarm
	NGroups                int
	PercentGroupIterations float64
}

func NewGroupBestSwarm(w float64, c1 float64, c2 float64, nParticles int, nGroups int, percentGroupIterations float64) *GroupBestSwarm {
	return &GroupBestSwarm{
		baseSwarm:              newBaseSwarm(w, c1, c2, nParticles),
		NGroups:                nGroups,
		PercentGroupIterations: percentGroupIterations,
	}
}
func (s *GroupBestSwarm) run(numIterations int, f func(left chan *ParticleStatus, right chan *ParticleStatus)) chan float64 {
	fitnessChan := make(chan float64, numIterations)
	go func() {
		numInGroup := s.config.NParticles / s.NGroups
		swarm := NewGlobalBestSwarm(s.config.W, s.config.C1, s.config.C2, s.NGroups)
		swarm.BindProblem(s.problem)
		for i := 0; i < s.NGroups; i++ {
			swarm.particles[i].Status = NewParticleStatusWithAnother(&s.particles[i*numInGroup].Status)
			swarm.particles[i].BestStatus = NewParticleStatusWithAnother(&s.particles[i*numInGroup].Status)
		}
		gBestStatus := NewParticleStatusWithAnother(&swarm.particles[0].Status)

		gg := s.NGroups - 1
		nGroupIterations := int(float64(numIterations) * s.PercentGroupIterations)
		for i := 0; i < numIterations; i++ {
			for g := 0; g < s.NGroups-1; g++ {
				for j := g * numInGroup; j < g*(numInGroup+1); j++ {
					swarm.particles[g].Status.OverwriteIfBetter(&s.particles[j].Status, s.better)
				}
			}
			for j := gg * numInGroup; j < s.config.NParticles; j++ {
				swarm.particles[gg].Status.OverwriteIfBetter(&s.particles[j].Status, s.better)
			}
			if nGroupIterations > 0 {
				_ = swarm.run(nGroupIterations, f)
				swarm.Wait()
			} else {
				for kk := 0; kk < s.NGroups; kk++ {
					swarm.particles[kk].BestStatus.Overwrite(&swarm.particles[kk].Status)
				}
			}
			for index := range swarm.particles {
				gBestStatus.OverwriteIfBetter(&swarm.particles[index].Status, s.better)
			}
			fitnessChan <- gBestStatus.Fitness
			for g := 0; g < s.NGroups-1; g++ {
				for j := g * numInGroup; j < g*(numInGroup+1); j++ {
					for k := 0; k < s.problem.Dims(); k++ {
						s.UpdateParticleVAndPosition(s.particles[j], k, swarm.particles[g].BestStatus.Positions[k])
					}
					s.particles[j].UpdateFitness()
				}
			}
			for j := gg * numInGroup; j < s.config.NParticles; j++ {
				for k := 0; k < s.problem.Dims(); k++ {
					s.UpdateParticleVAndPosition(s.particles[j], k, swarm.particles[gg].BestStatus.Positions[k])
				}
				s.particles[j].UpdateFitness()
			}
		}
		close(fitnessChan)
		s.stopCh <- true
	}()
	s.isRunning = true
	return fitnessChan
}
func (s *GroupBestSwarm) Run(numIterations int) chan float64 {
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
