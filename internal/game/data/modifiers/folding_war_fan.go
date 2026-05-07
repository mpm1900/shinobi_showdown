package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var foldingWarFanID = uuid.MustParse("c90c230c-0ae5-4b95-ab01-52a9cb5e0451")
var FoldingWarFan game.Modifier = game.Modifier{
	ID:          foldingWarFanID,
	GroupID:     &foldingWarFanID,
	Icon:        "folding_war_fan",
	Name:        "Folding War-Fan",
	Description: "Wind attacks deal 10% more damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&foldingWarFanID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureWind] += 0.1
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
