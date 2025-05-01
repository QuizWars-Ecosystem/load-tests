package scenario

import (
	"context"
	"fmt"
)

var GlobalAddr string

func SetGlobalAddr(addr string) {
	GlobalAddr = addr
}

func init() {
	registry = make(map[string]*Scenario)
}

type Scenario struct {
	Name      string
	Init      func() error
	Call      func(ctx context.Context) error
	DependsOn []string
}

var registry map[string]*Scenario

func Register(key string, s *Scenario) {
	registry[key] = s
}

func Get(key string) (*Scenario, bool) {
	s, ok := registry[key]
	return s, ok
}

func ResolveScenarioOrder(selected []string) ([]*Scenario, error) {
	visited := map[string]bool{}
	var result []*Scenario
	var visit func(name string) error

	visit = func(name string) error {
		if visited[name] {
			return nil
		}
		sc, ok := Get(name)
		if !ok {
			return fmt.Errorf("unknown scenario: %s", name)
		}
		for _, dep := range sc.DependsOn {
			if err := visit(dep); err != nil {
				return err
			}
		}
		visited[name] = true
		result = append(result, sc)
		return nil
	}

	for _, name := range selected {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}
