package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var dragonFlamePepperID = uuid.MustParse("d0676564-df6a-4716-a3c8-7ac4c5aaabff")
var DragonFlamePepper game.Modifier = game.Modifier{
	ID:          dragonFlamePepperID,
	GroupID:     &dragonFlamePepperID,
	Icon:        "dragon_flame_pepper",
	Name:        "Dragon-Flame Pepper",
	Description: "Fire attacks deal 10% more damage.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&dragonFlamePepperID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.NatureDamage[game.NatureFire] += 0.1
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
