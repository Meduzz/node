package node

import "github.com/Meduzz/helper/fp/slice"

type (
	Tripplet struct {
		Start  string `json:"start"`
		Action Action `json:"action"`
		End    string `json:"end"`
	}

	SimpleFlow[T any] struct {
		name         string
		participants map[string]Node[T]
		graph        []*Tripplet
		start        Node[T]
	}

	FlowBuilder[T any] interface {
		Relation(start Node[T], action Action, end Node[T])
	}
)

var (
	_ FlowBuilder[string] = &SimpleFlow[string]{}
	_ Node[string]        = &SimpleFlow[string]{}
)

func NewFlow[T any](name string, start Node[T], handler func(handler FlowBuilder[T])) Node[T] {
	dag := make(map[string]Node[T])

	dag[start.Name()] = start

	flow := &SimpleFlow[T]{
		name:         name,
		start:        start,
		participants: dag,
	}

	handler(flow)

	return flow
}

func (s *SimpleFlow[T]) Relation(start Node[T], action Action, end Node[T]) {
	t := &Tripplet{
		Start:  start.Name(),
		Action: action,
		End:    end.Name(),
	}

	s.graph = append(s.graph, t)
	s.participants[start.Name()] = start
	s.participants[end.Name()] = end
}

func (s *SimpleFlow[T]) Name() string {
	return s.name
}

func (s *SimpleFlow[T]) Exec(ctx T) (Action, error) {
	action, err := Execute(s.start, ctx)

	if err != nil {
		return action, err
	}

	next := s.next(s.start.Name(), action)

	if next == nil {
		return action, nil
	}

	for {
		action, err = Execute(next, ctx)

		if err != nil {
			return action, err
		}

		next := s.next(next.Name(), action)

		if next == nil {
			break
		}
	}

	return action, err
}

func (s *SimpleFlow[T]) next(current string, result Action) Node[T] {
	next := slice.Head(slice.Filter(s.graph, func(t *Tripplet) bool {
		return t.Start == current && t.Action == result
	}))

	if next == nil {
		return nil
	}

	return s.participants[next.End]
}
