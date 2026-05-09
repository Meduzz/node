package node_test

import (
	"fmt"
	"testing"

	"github.com/Meduzz/node"
)

type (
	mockNode struct {
		name      string
		execFunc  func(ctx string) (node.Action, error)
		callCount int
	}
)

var (
	_ node.Node[string] = (*mockNode)(nil)
)

func (m *mockNode) Exec(ctx string) (node.Action, error) {
	m.callCount++
	return m.execFunc(ctx)
}

func (m *mockNode) Name() string {
	return m.name
}

func TestIterator(t *testing.T) {
	t.Run("Successful iteration", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		results := []string{}
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				results = append(results, ctx)
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if action != node.Success {
			t.Errorf("expected Success, got %s", action)
		}
		if len(results) != 3 {
			t.Errorf("expected 3 results, got %d", len(results))
		}
		if results[0] != "a" || results[1] != "b" || results[2] != "c" {
			t.Errorf("unexpected results: %v", results)
		}
		if child.callCount != 3 {
			t.Errorf("expected 3 calls to child, got %d", child.callCount)
		}
	})

	t.Run("Non-Success action propagation", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				if ctx == "b" {
					return node.Retry, nil
				}
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if action != node.Retry {
			t.Errorf("expected Retry, got %s", action)
		}
		if child.callCount != 2 {
			t.Errorf("expected 2 calls to child (stopped at 'b'), got %d", child.callCount)
		}
	})

	t.Run("Error propagation", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		expectedErr := fmt.Errorf("something went wrong")
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				if ctx == "b" {
					return node.Error, expectedErr
				}
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if action != node.Error {
			t.Errorf("expected Error action, got %s", action)
		}
		if child.callCount != 2 {
			t.Errorf("expected 2 calls to child, got %d", child.callCount)
		}
	})

	t.Run("Custom success action", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		customSuccess := node.Action("custom_success")
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				return customSuccess, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, customSuccess)
		action, err := iterator.Exec(items)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if action != customSuccess {
			t.Errorf("expected %s, got %s", customSuccess, action)
		}
		if child.callCount != 3 {
			t.Errorf("expected 3 calls to child, got %d", child.callCount)
		}
	})

	t.Run("Non-success action with error", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		expectedErr := fmt.Errorf("some error")
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				if ctx == "b" {
					return node.Retry, expectedErr
				}
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if action != node.Retry {
			t.Errorf("expected Retry, got %s", action)
		}
		if child.callCount != 2 {
			t.Errorf("expected 2 calls to child, got %d", child.callCount)
		}
	})

	t.Run("Blank action with error", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		expectedErr := fmt.Errorf("blank action error")
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				if ctx == "b" {
					return "", expectedErr
				}
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if action != node.Error {
			t.Errorf("expected Error action, got %s", action)
		}
		if child.callCount != 2 {
			t.Errorf("expected 2 calls to child, got %d", child.callCount)
		}
	})

	t.Run("Success action with error", func(t *testing.T) {
		items := []string{"a", "b", "c"}
		fieldName := "items"
		itemKey := "item"

		expectedErr := fmt.Errorf("success with error")
		child := &mockNode{
			name: "child",
			execFunc: func(ctx string) (node.Action, error) {
				if ctx == "b" {
					return node.Success, expectedErr
				}
				return node.Success, nil
			},
		}

		iterator := node.NewIterator(fieldName, itemKey, child, node.Success)
		action, err := iterator.Exec(items)

		if err != expectedErr {
			t.Errorf("expected error %v, got %v", expectedErr, err)
		}
		if action != node.Error {
			t.Errorf("expected Error action, got %s", action)
		}
		if child.callCount != 2 {
			t.Errorf("expected 2 calls to child, got %d", child.callCount)
		}
	})
}
