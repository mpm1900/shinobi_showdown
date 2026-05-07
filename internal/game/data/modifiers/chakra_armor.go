package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var chakraArmorID = uuid.MustParse("9c1e932f-4af1-4543-851e-ae7444d7eebb")
var ChakraArmor = game.Modifier{
	ID:          chakraArmorID,
	GroupID:     &chakraArmorID,
	Name:        "Chakra Armor",
	Icon:        "chakra_armor",
	Show:        true,
	Description: "If at full health: take 50% less damage.",
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&chakraArmorID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter, game.FullHealthFilter),
			func(g game.Game, a game.Actor, ctx game.Context) game.Actor {
				a.DamageReduction[game.ChakraAttack] *= 0.5
				a.DamageReduction[game.Attack] *= 0.5
				return a
			},
		),
	},
}
