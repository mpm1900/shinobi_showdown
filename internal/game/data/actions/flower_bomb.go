package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var FlowerBomb = MakeFlowerBomb()

func MakeFlowerBomb() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Flower Bomb",
		Description: "30% chance to poison target.",
		Nature:      game.Ptr(game.NsWood),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(90),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("134b2304-7829-4739-864e-5e8b77bf0a41"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Push(modifiers.ChancePoison(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
