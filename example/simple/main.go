package main

import (
	"context"
	"github.com/lbuening/gaffa"
)

type app struct {
	gaffa.Implements[gaffa.Main]
	firstService gaffa.Ref[FirstService]
}

func run(ctx context.Context, a *app) error {
	a.Logger(ctx).Info("Running app")
	err := a.firstService.Get().DoSomething(ctx)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := gaffa.Run[app, *app](context.Background(), run)
	if err != nil {
		panic(err)
	}
}
