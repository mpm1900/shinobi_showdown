package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var LightningKunai = MakeLightningKunai()

func MakeLightningKunai() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Lightning Kunai",
		Description: "10% chance to paralyze target.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("5b5a4fe7-5c3c-4279-ae28-bfde927f8d8b"),
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
