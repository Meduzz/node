package node

func Execute(node Node, ctx map[string]any) (Action, error) {
	var err error
	lc, hasLc := node.(NodeLifecycleHooks)

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
