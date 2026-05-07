package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var focusedID = uuid.MustParse("007e83c6-282f-44bf-b5b1-37f6e0e1f734")

var Focused game.Modifier = game.Modifier{
	ID:          focusedID,
	GroupID:     &focusedID,
	Icon:        "focused",
	Name:        "Focused",
	Description: "Actions cannot be redirected.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&focusedID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.IgnoreRedirect = true
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
