package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var lightningModeID = uuid.MustParse("ca746d9d-53a2-4870-befd-95fe28108e1a")
var LightningMode = game.Modifier{
	ID:          lightningModeID,
	GroupID:     &lightningModeID,
	Name:        "Lightning Mode",
	Icon:        "electrified_speed",
	Description: "Doubled speed on Electrified Terrain.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&lightningModeID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Terrain != game.GameTerrainElectrified {
					return actor
				}

				actor.Stats[game.StatSpeed] = game.Round(float64(actor.Stats[game.StatSpeed]) * 2)
				return actor
			},
		),
	},
}
