package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var coercedID = uuid.MustParse("2da97c8a-85d1-41f6-baaf-29ac7967b12d")

func Coerced(actionID uuid.UUID) game.Modifier {
	return game.Modifier{
		ID:          uuid.New(),
		Name:        "Coerced",
		Description: "Coerced shinobi must used their last used action.",
		Icon:        "coerced",
		Show:        true,
		GroupID:     &coercedID,
		Duration:    5,
		ActorMutations: []game.ActorMutation{
			game.MakeActorMutation(
				&coercedID,
				game.MutPriorityDefault,
				game.ComposeAF(game.ActiveFilter, game.TargetFilter),
				func(g game.Game, a game.Actor, c game.Context) game.Actor {
					for i, _ := range a.Actions {
						a.Actions[i].State.Disabled = a.Actions[i].ID != actionID && !a.Actions[i].Meta.Switch
					}
					return a
				},
			),
		},
	}
}
