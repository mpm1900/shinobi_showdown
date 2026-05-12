package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var conductiveBracersID = uuid.MustParse("99266bb1-bb19-4a0d-81a7-5c2462b9184c")
var ConductiveBracers game.Modifier = game.Modifier{
	ID:          conductiveBracersID,
	GroupID:     &conductiveBracersID,
	Icon:        "conductive_bracers",
	Name:        "Conductive Bracers",
	Description: "Lightning attacks deal 1.2x damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&conductiveBracersID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureLightning] += 0.2
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
