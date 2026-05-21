package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var priorityFailureID = uuid.MustParse("7fdf20ca-b003-5082-8980-a0c6990169d0")
var PriorityFailure = game.Modifier{
	ID:          priorityFailureID,
	GroupID:     &priorityFailureID,
	Icon:        "priority_failure",
	Name:        "Priority Failure",
	Description: "Priority attacks are disabled.",
	Show:        true,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&priorityFailureID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				for i, action := range a.Actions {
					if !action.Config.Switch {
						if action.Priority > game.ActionPriorityDefault && action.Config.Power != nil {
							a.Actions[i].State.Disabled = true
						}
					}
				}

				return a
			},
		),
	},
}
