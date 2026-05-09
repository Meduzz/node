package node

type Iterator[T any] struct {
	fieldName     string
	itemKey       string
	child         Node[T]
	successAction Action
}

func NewIterator[T any](fieldName, itemKey string, child Node[T], successAction Action) *Iterator[T] {
	return &Iterator[T]{
		fieldName:     fieldName,
		itemKey:       itemKey,
		child:         child,
		successAction: successAction,
	}
}

func (i *Iterator[T]) Name() string {
	return "iterator"
}

func (i *Iterator[T]) Exec(ctx []T) (Action, error) {

	for idx := range ctx {
		item := ctx[idx]

		action, err := Execute(i.child, item)

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
