package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GreatTreeSpear = MakeGreatTreeSpear()

func MakeGreatTreeSpear() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Great Tree Spear",
		Description: "30% chance to poison target.",
		Nature:      game.Ptr(game.NsWood),
		Accuracy:    game.Ptr(70),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(90),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("00888b4b-973f-5bf5-9a41-bba1c9b629b8"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChancePoison(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
