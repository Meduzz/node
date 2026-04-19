package node

import (
	"fmt"
	"reflect"
)

type Iterator struct {
	fieldName     string
	itemKey       string
	child         Node
	successAction Action
}

func NewIterator(fieldName, itemKey string, child Node, successAction Action) *Iterator {
	return &Iterator{
		fieldName:     fieldName,
		itemKey:       itemKey,
		child:         child,
		successAction: successAction,
	}
}

func (i *Iterator) Name() string {
	return "iterator"
}

func (i *Iterator) Exec(ctx map[string]any) (Action, error) {
	val, ok := ctx[i.fieldName]

	if !ok {
		return Error, fmt.Errorf("field %s not found in context", i.fieldName)
	}

	v := reflect.ValueOf(val)

	if v.Kind() != reflect.Slice {
		return Error, fmt.Errorf("field %s is not a slice", i.fieldName)
	}

	for idx := 0; idx < v.Len(); idx++ {
		item := v.Index(idx).Interface()
		ctx[i.itemKey] = item

		action, err := Execute(i.child, ctx)
		if action != i.successAction {
			if action == "" {
				return Error, err
			}
			return action, err
		}

		if err != nil {
			return Error, err
		}
	}

	return i.successAction, nil
}
