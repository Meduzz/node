package node

import "github.com/Meduzz/helper/fp/slice"

type (
	Tripplet struct {
		Start  string `json:"start"`
		Action Action `json:"action"`
		End    string `json:"end"`
	}

	SimpleFlow struct {
		name         string
		participants map[string]Node
		graph        []*Tripplet
		start        Node
	}

	FlowBuilder interface {
		Relation(start Node, action Action, end Node)
	}
)

var (
	_ FlowBuilder = &SimpleFlow{}
	_ Node        = &SimpleFlow{}
)

func NewFlow(name string, start Node, handler func(handler FlowBuilder)) Node {
	dag := make(map[string]Node)

	dag[start.Name()] = start

	flow := &SimpleFlow{
		name:         name,
		start:        start,
		participants: dag,
	}

	handler(flow)

	return flow
}

func (s *SimpleFlow) Relation(start Node, action Action, end Node) {
	t := &Tripplet{
		Start:  start.Name(),
		Action: action,
		End:    end.Name(),
	}

	s.graph = append(s.graph, t)
	s.participants[start.Name()] = start
	s.participants[end.Name()] = end
}

func (s *SimpleFlow) Name() string {
	return s.name
}

func (s *SimpleFlow) Exec(ctx map[string]any) (Action, error) {
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

func (s *SimpleFlow) next(current string, result Action) Node {
	next := slice.Head(slice.Filter(s.graph, func(t *Tripplet) bool {
		return t.Start == current && t.Action == result
	}))

	if next == nil {
		return nil
	}

	return s.participants[next.End]
}
