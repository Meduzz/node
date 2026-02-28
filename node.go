package node

type (
	Action string

	Settings struct {
		Retries int `json:"retries,omitempty"`
		// TODO timeout?
	}

	Node interface {
		Name() string
		Exec(map[string]any) (Action, error)
	}

	NodeLifecycleHooks interface {
		Pre(map[string]any) error
		Post(map[string]any, Action, error) (Action, error)
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
