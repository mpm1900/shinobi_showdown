package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var LightningHound = MakeLightningHound()

func MakeLightningHound() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Lightning Hound",
		Description: "10% chance to paralyze target.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("5a35a0b0-160a-4b73-9b81-bb301e7c8f7e"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Push(modifiers.ChanceParalysis(action_config, g, context, target, 10))
			}

			return transactions.Build()
		},
	})
}
