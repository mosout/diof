package main

import (
	"fmt"
	"github.com/mosout/diof/pso"
	"github.com/mosout/diof/pso/problems"
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	problem := problems.NewRAS(10, 5.0)
	s := pso.NewGlobalBestSwarm(0.5, 2, 2, 100)
	s.BindProblem(problem)
	err := s.Connect("localhost", "7365", pso.ServerParams{
		Namespace: "test",
		Target:    problem.Target(),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := s.Disconnect(); err != nil {
			log.Fatal(err)
		}
	}()
	fitnessChan := s.Run(500)
	for fitness := range fitnessChan {
		fmt.Println(fitness)
	}
	defer s.Wait()
}
