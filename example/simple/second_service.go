package main

import (
	"context"
	"github.com/lbuening/gaffa"
)

type SecondService interface {
	DoSomethingElse(ctx context.Context)
}

type secondServiceFirstImpl struct {
	gaffa.Implements[SecondService] `bidde:"default"`
}

func (s *secondServiceFirstImpl) DoSomethingElse(ctx context.Context) {
	s.Logger(ctx).Info("Doing something else on default")
}

type secondServiceSecondImpl struct {
	gaffa.Implements[SecondService] `bidde:"cloud"`
}

func (s *secondServiceSecondImpl) DoSomethingElse(ctx context.Context) {
	s.Logger(ctx).Info("Doing something else on cloud")
}
