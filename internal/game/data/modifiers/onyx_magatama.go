package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var onyxMagatamaID = uuid.MustParse("84369b0e-b364-4a83-b4d2-979477ec50f5")
var OnyxMagatama game.Modifier = game.Modifier{
	ID:          onyxMagatamaID,
	GroupID:     &onyxMagatamaID,
	Icon:        "onyx_magatama",
	Name:        "Onyx Magatama",
	Description: "Yin attacks deal 1.2x damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&onyxMagatamaID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureYin] += 0.2
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
