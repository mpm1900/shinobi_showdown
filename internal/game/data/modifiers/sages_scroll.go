package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var sagesScrollID = uuid.MustParse("0a88ef04-e105-47ec-ae74-3bec50bd62cd")
var SagesScroll game.Modifier = game.Modifier{
	ID:          sagesScrollID,
	GroupID:     &sagesScrollID,
	Icon:        "sages_scroll",
	Name:        "Sage's Scroll",
	Description: "Yang attacks deal 10% more damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&sagesScrollID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureYang] += 0.1
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
