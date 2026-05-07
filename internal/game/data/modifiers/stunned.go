package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var stunnedID = uuid.MustParse("3f2d768d-3ed7-5b53-8655-d738e2ca3a11")
var Stunned = game.Modifier{
	ID:          stunnedID,
	Name:        "Stunned",
	Description: "Stunned: cannot act.",
	Icon:        "stunned",
	Show:        true,
	GroupID:     &stunnedID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&stunnedID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.TargetFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.Stunned = true
				return a
			},
		),
	},
}
