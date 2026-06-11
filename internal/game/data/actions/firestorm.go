package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Firestorm = MakeFirestorm()

func MakeFirestorm() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Firestorm",
		Description: "10% chance to burn target. Lowers user's Chakra attack by 2 stages.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(100),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("5756b76d-dd39-460c-b5fa-431b80200f3b"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceBurn(action_config, g, context, target, 10))
			}

			mod := modifiers.ChakraAttackDown2Source
			mut := mutations.AddModifiers(false, mod)
			transactions.Push(game.MakeTransaction(mut, context))

			return transactions.Build()
		},
	})
}
