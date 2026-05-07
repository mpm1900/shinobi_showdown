package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var protectedID = uuid.MustParse("753f21af-6b67-50ff-a386-40528a73e62e")
var Protected = game.Modifier{
	ID:          uuid.MustParse("fa58e0c1-01aa-58ef-b42d-0af89ee4555c"),
	Name:        "Protected",
	Description: "Protected from damage and modifiers.",
	Icon:        "protected",
	Show:        true,
	GroupID:     &protectedID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&protectedID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.Protected = true
				a.Safeguarded = true
				return a
			},
		),
	},
}

var spreadProtectedID = uuid.MustParse("4a4f0fc5-96a6-4b6a-9e8e-664964fab220")
var SpreadProtected = game.Modifier{
	ID:          uuid.MustParse("fa58e0c1-01aa-58ef-b42d-0af89ee4555c"),
	Name:        "Spread Protected",
	Description: "Protected from damage and modifiers that target more than one shinobi.",
	Icon:        "protected",
	Show:        true,
	GroupID:     &spreadProtectedID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&spreadProtectedID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				atx := g.ActiveTransaction
				if atx == nil {
					return a
				}

				targets := g.GetTargets(atx.Context)
				if len(targets) <= 1 {
					return a
				}

				a.Protected = true
				a.Safeguarded = true
				return a
			},
		),
	},
}
