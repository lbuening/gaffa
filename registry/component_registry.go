package registry

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"sync"
)

var globalRegistry = componentRegistry{
	registrations: []*Registration{},
	initialized:   map[string]any{},
}

type Registration struct {
	Iface reflect.Type
	Impl  reflect.Type
}

type componentRegistry struct {
	m             sync.Mutex
	registrations []*Registration
	initialized   map[string]any
}

func (r *componentRegistry) Register(reg Registration) error {
	r.m.Lock()
	defer r.m.Unlock()
	r.registrations = append(r.registrations, &reg)
	return nil
}

func (r *componentRegistry) findCandidatesForIface(t reflect.Type) (map[string]*Registration, error) {
	r.m.Lock()
	defer r.m.Unlock()
	var interfaceCandidates []*Registration
	for _, reg := range r.registrations {
		if reg.Iface != t {
			continue
		}
		interfaceCandidates = append(interfaceCandidates, reg)
	}

	candidates := map[string]*Registration{}
	for _, candidate := range interfaceCandidates {
		tag, ok := extractTagFromType(candidate.Impl)
		if !ok && len(interfaceCandidates) == 1 {
			return map[string]*Registration{"default": candidate}, nil
		}
		if !ok {
			return nil, fmt.Errorf("multiple implementations of %v, please specify one with a gaffa tag", t)
		}
		if _, ok := candidates[tag]; ok {
			return nil, fmt.Errorf("duplicate gaffa tag %q for %v", tag, t)
		}
		candidates[tag] = candidate
	}

	return candidates, nil
}

func (r *componentRegistry) findCandidatesForImpl(t reflect.Type) (map[string]*Registration, error) {
	r.m.Lock()
	defer r.m.Unlock()
	var implCandidates []*Registration
	for _, reg := range r.registrations {
		if reg.Impl != t {
			continue
		}
		implCandidates = append(implCandidates, reg)
	}

	if len(implCandidates) == 1 {
		return map[string]*Registration{"default": implCandidates[0]}, nil
	}

	candidates := map[string]*Registration{}
	// Check for a candidate with the given tag.
	for _, candidate := range implCandidates {
		tag, ok := extractTagFromType(candidate.Impl)
		if !ok {
			return nil, fmt.Errorf("multiple implementations of %v, please specify one with a gaffa tag", t)
		}
		if _, ok := candidates[tag]; ok {
			return nil, fmt.Errorf("duplicate gaffa tag %q for %v", tag, t)
		}
		candidates[tag] = candidate
	}

	return candidates, nil
}

func extractTagFromType(t reflect.Type) (string, bool) {
	if t.Kind() != reflect.Struct {
		return "", false
	}
	if t.NumField() == 0 {
		return "", false
	}
	f := t.Field(0)
	return f.Tag.Lookup("gaffa")
}

func (r *componentRegistry) initialize(reg *Registration, tag string) (any, error) {
	v := reflect.New(reg.Impl)
	obj := v.Interface()

	err := SetLogger(obj, slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err != nil {
		return nil, err
	}

	if HasRefs(obj) {
		err := FillRefs(obj, func(t reflect.Type) (any, error) {
			return GetIntf(t, tag)
		})
		if err != nil {
			return nil, err
		}
	}

	r.initialized[identifier(reg)] = obj
	return obj, nil
}

func identifier(reg *Registration) string {
	h := sha1.New()
	h.Write([]byte(reg.Iface.PkgPath()))
	h.Write([]byte(reg.Iface.String()))
	h.Write([]byte(reg.Impl.PkgPath()))
	h.Write([]byte(reg.Impl.String()))
	return hex.EncodeToString(h.Sum(nil))
}

func Register(reg Registration) {
	err := globalRegistry.Register(reg)
	if err != nil {
		panic(err)
	}
}

func GetImpl(r reflect.Type, tag string) (any, error) {
	if tag == "" {
		tag = "default"
	}
	candidates, err := globalRegistry.findCandidatesForImpl(r)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 1 {
		return globalRegistry.initialize(candidates["default"], tag)
	}
	reg, ok := candidates[tag]
	if !ok {
		if tag == "default" {
			return nil, fmt.Errorf("component interface %v not found; maybe you forgot to run gaffa generate", r)
		}
		reg, ok = candidates["default"]
		if !ok {
			return nil, fmt.Errorf("component interface %v not found; maybe you forgot to run gaffa generate", r)
		}
	}
	return globalRegistry.initialize(reg, tag)
}

func GetIntf(r reflect.Type, tag string) (any, error) {
	if tag == "" {
		tag = "default"
	}
	candidates, err := globalRegistry.findCandidatesForIface(r)
	if err != nil {
		return nil, err
	}
	reg, ok := candidates[tag]
	if !ok {
		if tag == "default" {
			return nil, fmt.Errorf("component interface %v not found; maybe you forgot to run gaffa generate", r)
		}
		reg, ok = candidates["default"]
		if !ok {
			return nil, fmt.Errorf("component interface %v not found; maybe you forgot to run gaffa generate", r)
		}
	}
	return globalRegistry.initialize(reg, tag)
}
