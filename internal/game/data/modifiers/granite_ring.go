package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var graniteRingID = uuid.MustParse("9de939ad-24eb-46cd-b0a8-536fccf88234")
var GraniteRing game.Modifier = game.Modifier{
	ID:          graniteRingID,
	GroupID:     &graniteRingID,
	Icon:        "granite_ring",
	Name:        "Granite Ring",
	Description: "Earth attacks deal 1.2x damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&graniteRingID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureEarth] += 0.2
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
