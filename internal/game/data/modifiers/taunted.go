package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var tauntedID = uuid.MustParse("a1d1b733-5c33-50ff-a20a-b090ee150650")
var Taunted = game.Modifier{
	ID:          uuid.MustParse("2ae07cad-2c15-540d-82f7-91308b3c32d0"),
	Name:        "Taunted",
	Icon:        "taunted",
	Description: "Attacking actions only.",
	Show:        true,
	GroupID:     &tauntedID,
	Duration:    4,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&tauntedID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.TargetFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				for i := range a.Actions {
					a.Actions[i].Disabled = !a.Actions[i].Meta.Switch && a.Actions[i].Config.Power == nil
				}

				return a
			},
		),
	},
}
