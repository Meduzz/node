package node

type (
	Action string

	Settings struct {
		Retries int `json:"retries,omitempty"`
		// TODO timeout?
	}

	Node[T any] interface {
		Name() string
		Exec(ctx T) (Action, error)
	}

	NodeLifecycleHooks[T any] interface {
		Pre(ctx T) error
		Post(ctx T, action Action, err error) (Action, error)
	}

	NodeSettingsHooks interface {
		Settings() *Settings
	}
)

const (
	Success = Action("success")
	Retry   = Action("retry")
	Error   = Action("error")
)
