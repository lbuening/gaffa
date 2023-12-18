package gaffa

import (
	"context"
	"github.com/lbuening/gaffa/registry"
	"log/slog"
	"reflect"
)

// Implements is a type that is being embedded inside a component
// implementation struct to indicate that the struct implements a component of
// type T. For example, consider a Cache component.
//
//	type Cache interface {
//	    Get(ctx context.Context, key string) (string, error)
//	    Put(ctx context.Context, key, value string) error
//	}
//
// A concrete type that implements the Cache component is written as follows:
//
//	type lruCache struct {
//	    weaver.Implements[Cache]
//	    ...
//	}
//
// Because Implements is embedded inside the component implementation, methods
// of Implements are available as methods of the component implementation type
// and can be invoked directly. For example, given an instance c of type
// lruCache, we can call c.Logger().
type Implements[T any] struct {
	// Component logger.
	logger *slog.Logger

	// Given a component implementation type, there is currently no nice way,
	// using reflection, to get the corresponding component interface type [1].
	// The component_interface_type field exists to make it possible.
	//
	// [1]: https://github.com/golang/go/issues/54393.
	//
	//lint:ignore U1000 See comment above.
	component_interface_type T

	// We embed implementsImpl so that component implementation structs
	// implement the Unrouted interface by default but implement the
	// RoutedBy[T] interface when they embed WithRouter[T].
	implementsImpl
}

// Logger returns a logger that associates its log entries with this component.
func (i *Implements[T]) Logger(_ context.Context) *slog.Logger {
	return i.logger
}

func (i *Implements[T]) setLogger(logger *slog.Logger) {
	i.logger = logger
}

// implements is a method that can only be implemented inside the weaver
// package. It exists so that a component struct that embeds Implements[T]
// implements the InstanceOf[T] interface.
//
//lint:ignore U1000 implements is used by InstanceOf.
func (i *Implements[T]) implements(T) {}

// InstanceOf is the interface implemented by a struct that embeds
// gaffa.Implements[T].
type InstanceOf[T any] interface {
	implements(T)
}

// See Implements.implementsImpl.
type implementsImpl struct{}

// Ref is a field that can be placed inside a component implementation
// struct. T must be a component type. Service Weaver will automatically
// fill such a field with a handle to the corresponding component.
type Ref[T any] struct {
	value T
}

// Get returns a handle to the component of type T.
func (r Ref[T]) Get() T { return r.value }

// isRef is an internal method that is only implemented by Ref[T] and is
// used internally to check that a value is of type Ref[T].
func (r Ref[T]) isRef() {}

// setRef sets the underlying value of a Ref.
func (r *Ref[T]) setRef(value any) {
	r.value = value.(T)
}

// Main is the interface implemented by an application's main component.
type Main interface{}

// PointerToMain is a type constraint that asserts *T is an instance of Main
// (i.e. T is a struct that embeds weaver.Implements[weaver.Main]).
type PointerToMain[T any] interface {
	*T
	InstanceOf[Main]
}

func Run[T any, _ PointerToMain[T]](ctx context.Context, app func(context.Context, *T) error) error {
	main, err := registry.GetImpl(Type[T](), "default")
	if err != nil {
		return err
	}
	return app(ctx, main.(*T))
}

func Type[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
