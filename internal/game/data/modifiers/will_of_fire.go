package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var willOfFireID = uuid.MustParse("04761c26-1a29-42ee-8e8f-532ffa1a889e")
var WillOfFire game.Modifier = game.Modifier{
	ID:                willOfFireID,
	GroupID:           &willOfFireID,
	Icon:              "will_of_fire",
	Name:              "Will Of Fire",
	Description:       "Fire attacks deal 50% more damage.",
	ParentDescription: "At below 50% HP: Fire attacks deal 50% more damage.",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&willOfFireID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.IsAtOrBelowHealthRatio(0.5), game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureFire] += 0.5
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
