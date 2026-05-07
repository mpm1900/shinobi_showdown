package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var disabledID = uuid.MustParse("2da97c8a-85d1-41f6-baaf-29ac7967b12d")

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
