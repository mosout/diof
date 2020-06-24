package rpc

import (
	"context"
	"github.com/mosout/diof/models"
	"github.com/mosout/diof/pso"
	"github.com/smallnest/rpcx/share"
	"sync"
)

type PSO struct {
	bestStatus        sync.Map
	connectionCounter sync.Map
	better            sync.Map
}

func (s *PSO) Enter(ctx context.Context, args *models.Empty, reply *models.Empty) error {
	reqMeta := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	namespace, target, uuid := reqMeta["namespace"], reqMeta["target"], reqMeta["id"]
	if _, ok := s.better.Load(namespace); !ok {
		s.better.Store(namespace, pso.SelectBetter(target))
	}
	if conns, ok := s.connectionCounter.Load(namespace); !ok {
		s.connectionCounter.Store(namespace, map[string]bool{
			uuid: true,
		})
	} else {
		connsMap := conns.(map[string]bool)
		connsMap[uuid] = true
		s.connectionCounter.Store(namespace, connsMap)
	}
	return nil
}
func (s *PSO) Exit(ctx context.Context, args *models.Empty, reply *models.Empty) error {
	reqMeta := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	namespace, uuid := reqMeta["namespace"], reqMeta["id"]
	if conns, ok := s.connectionCounter.Load(namespace); ok {
		connsMap := conns.(map[string]bool)
		delete(connsMap, uuid)
		if len(connsMap) == 0 {
			s.connectionCounter.Delete(namespace)
			s.bestStatus.Delete(namespace)
			s.better.Delete(namespace)
		}
		s.connectionCounter.Store(namespace, connsMap)
	}
	return nil
}
func (s *PSO) Push(ctx context.Context, ps *pso.ParticleStatus, reply *models.Empty) error {
	reqMeta := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	namespace := reqMeta["namespace"]
	if old, ok := s.bestStatus.Load(namespace); !ok {
		s.bestStatus.Store(namespace, ps)
	} else {
		oldStatus := old.(*pso.ParticleStatus)
		better, _ := s.better.Load(namespace)
		betterFunc := better.(func(float64, float64) bool)
		if betterFunc(ps.Fitness, oldStatus.Fitness) {
			s.bestStatus.Store(namespace, ps)
		}
	}
	return nil
}
func (s *PSO) Pull(ctx context.Context, args *models.Empty, ps *pso.ParticleStatus) error {
	reqMeta := ctx.Value(share.ReqMetaDataKey).(map[string]string)
	namespace := reqMeta["namespace"]
	if data, ok := s.bestStatus.Load(namespace); ok {
		ps.Overwrite(data.(*pso.ParticleStatus))
	}
	return nil
}
