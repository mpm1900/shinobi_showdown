package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var purePowerID = uuid.MustParse("7379e579-82b9-40d1-8fcd-5afedde9e739")
var PurePower = game.Modifier{
	ID:          purePowerID,
	GroupID:     &purePowerID,
	Name:        "Pure Power",
	Icon:        "attack_double",
	Description: "This shinobi's Attack stat is doubled.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&purePowerID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stats[game.StatAttack] = game.Round(float64(actor.Stats[game.StatAttack]) * 2)
				return actor
			},
		),
	},
}
