package node_test

import (
	"fmt"
	"testing"

	"github.com/Meduzz/node"
)

type (
	sprinter struct {
		format string
		input  string
		target string
	}
)

var (
	_     node.Node = &sprinter{}
	INPUT           = "input"
	HELLO           = "hello"
	BYE             = "bye"
)

func TestFlows(t *testing.T) {
	start := &sprinter{"Hello %v!", INPUT, HELLO}
	end := &sprinter{"Bye cruel %v!", INPUT, BYE}
	printer := node.NewFlow("test", start, func(handler node.FlowBuilder) {
		handler.Relation(start, node.Success, end)
	})
	ctx := make(map[string]any)
	ctx[INPUT] = "world"

	action, err := printer.Exec(ctx)

	if err != nil {
		t.Error(err)
	}

	if action != node.Success {
		t.Errorf("Action was not %s but %s", node.Success, action)
	}

	result, exists := ctx[HELLO]

	if !exists {
		t.Error("HELLO was not set")
	}

	if result != "Hello world!" {
		t.Errorf("result was not the expected: '%v'", result)
	}

	result, exists = ctx[BYE]

	if !exists {
		t.Error("BYE was not set")
	}

	if result != "Bye cruel world!" {
		t.Errorf("result was not the expected: '%v'", result)
	}
}

func (s *sprinter) Exec(ctx map[string]any) (node.Action, error) {
	anyInput, exists := ctx[s.input]

	if !exists {
		return node.Error, fmt.Errorf("%s was not set", s.input)
	}

	ctx[s.target] = fmt.Sprintf(s.format, anyInput)
	return node.Success, nil
}

func (s *sprinter) Name() string {
	return s.target
}
