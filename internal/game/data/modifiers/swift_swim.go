package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var swiftSwimID = uuid.MustParse("15427a62-2078-400f-ad55-7f3d6b85f5ca")
var SwiftSwim = game.Modifier{
	ID:          swiftSwimID,
	GroupID:     &swiftSwimID,
	Name:        "Swift Swim",
	Icon:        "rain_speed",
	Description: "Doubled speed while raining.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&swiftSwimID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				state, _ := g.GetState(context)
				if state.Weather != game.GameWeatherRain {
					return actor
				}

				actor.Stats[game.StatSpeed] = game.Round(float64(actor.Stats[game.StatSpeed]) * 2)
				return actor
			},
		),
	},
}
