package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var flyingID = uuid.MustParse("2b7a0e44-ed1f-4d0e-96e5-998c9f4e5272")
var Flying = game.Modifier{
	ID:          flyingID,
	GroupID:     &flyingID,
	Icon:        "flying",
	Name:        "Flying",
	Description: "(\"grounded\" is default)",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&flyingID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter),
			func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
				a.State = game.ActorStateFlying
				return a
			},
		),
	},
}
