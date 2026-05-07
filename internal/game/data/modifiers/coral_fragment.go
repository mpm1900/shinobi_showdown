package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var coralFragmentID = uuid.MustParse("e1aafdaf-eb75-4184-a285-407802449e04")
var CoralFragment game.Modifier = game.Modifier{
	ID:          coralFragmentID,
	GroupID:     &coralFragmentID,
	Icon:        "coral_fragment",
	Name:        "Coral Fragment",
	Description: "Water attacks deal 10% more damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&coralFragmentID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureWater] += 0.1
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
