package main

import (
	"context"
	"github.com/lbuening/gaffa"
)

type FirstService interface {
	DoSomething(context.Context) error
}

type firstServiceImpl struct {
	gaffa.Implements[FirstService]
	secondService gaffa.Ref[SecondService]
}

func (f *firstServiceImpl) DoSomething(ctx context.Context) error {
	f.Logger(ctx).Info("Doing something")
	f.secondService.Get().DoSomethingElse(ctx)
	return nil
}
