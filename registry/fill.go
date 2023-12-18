package registry

import (
	"log/slog"
	"reflect"
)

var (
	SetLogger func(impl any, logger *slog.Logger) error

	HasRefs func(impl any) bool

	FillRefs func(impl any, get func(reflect.Type) (any, error)) error
)
