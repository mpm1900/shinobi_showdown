package game

type TurnPhase string

const (
	TurnInit    TurnPhase = "init"
	TurnStart   TurnPhase = "start"
	TurnMain    TurnPhase = "main"
	TurnEnd     TurnPhase = "end"
	TurnCleanup TurnPhase = "cleanup"
)

type ActionStatus string

const (
	ActionStatusNone   ActionStatus = "none"
	ActionStatusPre    ActionStatus = "pre"
	ActionStatusActive ActionStatus = "active"
	ActionStatusPost   ActionStatus = "post"
)

type Turn struct {
	Count        int          `json:"count"`
	ActionStatus ActionStatus `json:"-"`
	Phase        TurnPhase    `json:"phase"`
}

func NewTurn() Turn {
	return Turn{
		Count:        0,
		Phase:        TurnInit,
		ActionStatus: ActionStatusNone,
	}
}
