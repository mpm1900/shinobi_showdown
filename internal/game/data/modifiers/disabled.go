package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var disabledID = uuid.MustParse("cea0796c-df52-466f-8474-9dc06ec9db6f")

func Disabled(actionID uuid.UUID) game.Modifier {
	return game.Modifier{
		ID:          uuid.New(),
		Name:        "Disabled",
		Description: "Disabled shinobi cannot use specific actions.",
		Icon:        "disabled",
		Show:        true,
		GroupID:     &disabledID,
		Duration:    5,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				&disabledID,
				game.MutPriorityDefault,
				game.ComposeAF(game.ActiveFilter, game.TargetFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					for i, _ := range a.Actions {
						if a.Actions[i].ID == actionID {
							a.Actions[i].Disabled = true
						}
					}
					return a
				},
			),
		},
	}
}
