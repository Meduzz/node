package node

func Execute[T any](node Node[T], ctx T) (Action, error) {
	var err error
	lc, hasLc := node.(NodeLifecycleHooks[T])

	if hasLc {
		err = lc.Pre(ctx)

		if err != nil {
			return lc.Post(ctx, Error, err)
		}
	}

	action, err := node.Exec(ctx)

	if hasLc {
		return lc.Post(ctx, action, err)
	}

	return action, err
}
