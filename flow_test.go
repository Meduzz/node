package node_test

import (
	"fmt"
	"testing"

	"github.com/Meduzz/node"
)

type (
	sprinter struct {
		format string
		target string
		value  string
	}
)

var (
	_     node.Node[string] = &sprinter{}
	HELLO                   = "hello"
	BYE                     = "bye"
)

func TestFlows(t *testing.T) {
	start := &sprinter{"Hello %v!", HELLO, ""}
	end := &sprinter{"Bye cruel %v!", BYE, ""}
	printer := node.NewFlow("test", start, func(handler node.FlowBuilder[string]) {
		handler.Relation(start, node.Success, end)
	})

	ctx := "world"

	action, err := printer.Exec(ctx)

	if err != nil {
		t.Error(err)
	}

	if action != node.Success {
		t.Errorf("Action was not %s but %s", node.Success, action)
	}

	if start.value != "Hello world!" {
		t.Errorf("result was not the expected: '%v'", start.value)
	}

	if end.value != "Bye cruel world!" {
		t.Errorf("result was not the expected: '%v'", end.value)
	}
}

func (s *sprinter) Exec(ctx string) (node.Action, error) {
	s.value = fmt.Sprintf(s.format, ctx)
	return node.Success, nil
}

func (s *sprinter) Name() string {
	return s.target
}
