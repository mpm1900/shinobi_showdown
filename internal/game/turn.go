package game

type TurnPhase string

const (
	TurnInit    TurnPhase = "init"
	TurnStart   TurnPhase = "start"
	TurnMain    TurnPhase = "main"
	TurnEnd     TurnPhase = "end"
	TurnCleanup TurnPhase = "cleanup"
)

type Turn struct {
	Count     int       `json:"count"`
	PreAction bool      `json:"-"`
	Phase     TurnPhase `json:"phase"`
}

func NewTurn() Turn {
	return Turn{
		Count:     0,
		Phase:     TurnInit,
		PreAction: false,
	}
}
