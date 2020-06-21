package pso

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mosout/diof/models"
	"github.com/mosout/diof/pso/problems"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
	"math/rand"
)

type baseSwarm struct {
	config    swarmConfig
	particles []*Particle
	vBounds   [][2]float64
	problem   problems.IProblem
	better    func(float64, float64) bool
	stopCh    chan bool
	isRunning bool
	rpcClient *client.Client
	uuid      string
	namespace string
}

func newBaseSwarm(w float64, c1 float64, c2 float64, nParticles int) baseSwarm {
	s := baseSwarm{
		config: swarmConfig{
			W:          w,
			C1:         c1,
			C2:         c2,
			NParticles: nParticles,
		},
		stopCh: make(chan bool),
	}
	s.particles = make([]*Particle, nParticles)
	return s
}
func (s *baseSwarm) BindProblem(p problems.IProblem) {
	s.problem = p
	s.better = SelectBetter(p.Target())
	s.vBounds = make([][2]float64, p.Dims())
	lowBounds := make([]float64, p.Dims())
	intervals := make([]float64, p.Dims())
	bounds := p.Bounds()
	for i := 0; i < p.Dims(); i++ {
		lowBounds[i] = bounds[i][0]
		intervals[i] = bounds[i][1] - bounds[i][0]
		vMax := intervals[i] * 0.2
		s.vBounds[i] = [2]float64{-vMax, vMax}
	}
	for i := 0; i < s.config.NParticles; i++ {
		s.particles[i] = NewParticle(lowBounds, intervals, p.Dims(), s)
	}
}
func (s *baseSwarm) DealOutOfBound(p *Particle, dim int) {
	if p.V[dim] > s.vBounds[dim][1] {
		p.V[dim] = s.vBounds[dim][1]
	}
	if p.V[dim] < s.vBounds[dim][0] {
		p.V[dim] = s.vBounds[dim][0]
	}
	if p.Status.Positions[dim] > s.problem.Bounds()[dim][1] {
		p.Status.Positions[dim] = s.problem.Bounds()[dim][1]
	}
	if p.Status.Positions[dim] < s.problem.Bounds()[dim][0] {
		p.Status.Positions[dim] = s.problem.Bounds()[dim][0]
	}
}
func (s *baseSwarm) UpdateParticleVAndPosition(p *Particle, dim int, bestPosition float64) {
	p.V[dim] = s.config.W*p.V[dim] +
		s.config.C1*rand.Float64()*(bestPosition-p.Status.Positions[dim]) +
		s.config.C2*rand.Float64()*(p.BestStatus.Positions[dim]-p.Status.Positions[dim])
	p.Status.Positions[dim] += p.V[dim]
	s.DealOutOfBound(p, dim)
}
func (s *baseSwarm) Wait() {
	if s.isRunning {
		<-s.stopCh
		s.isRunning = false
	}
}
func (s *baseSwarm) Connect(address string, port string, config ServerParams) error {
	s.rpcClient = client.NewClient(client.DefaultOption)
	if err := s.rpcClient.Connect("tcp", fmt.Sprintf("%s:%s", address, port)); err != nil {
		return err
	}
	s.uuid = uuid.New().String()
	s.namespace = config.Namespace
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{
		"id":        s.uuid,
		"namespace": config.Namespace,
		"target":    config.Target,
	})
	if err := s.rpcClient.Call(ctx, "PSO", "Enter", &models.Empty{}, &models.Empty{}); err != nil {
		return err
	}
	return nil
}
func (s *baseSwarm) Disconnect() error {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{
		"id":        s.uuid,
		"namespace": s.namespace,
	})
	if err := s.rpcClient.Call(ctx, "PSO", "Exit", &models.Empty{}, &models.Empty{}); err != nil {
		return err
	}
	if err := s.rpcClient.Close(); err != nil {
		return err
	}
	s.rpcClient = nil
	s.uuid = ""
	s.namespace = ""
	return nil
}
func (s *baseSwarm) Push(ps *ParticleStatus) {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{
		"namespace": s.namespace,
	})
	_ = s.rpcClient.Call(ctx, "PSO", "Push", ps, &models.Empty{})
}
func (s *baseSwarm) Pull() *ParticleStatus {
	ctx := context.WithValue(context.Background(), share.ReqMetaDataKey, map[string]string{
		"namespace": s.namespace,
	})
	ps := NewParticleStatus(s.problem.Dims())
	_ = s.rpcClient.Call(ctx, "PSO", "Pull", &models.Empty{}, &ps)
	return &ps
}
