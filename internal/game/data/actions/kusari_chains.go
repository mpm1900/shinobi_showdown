package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var KusariChains = MakeKusariChains()

func MakeKusariChains() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Kusari Chains",
		Description: "30% chance to stun target.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Bukijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("fba5ad6e-d2ee-4b3a-b524-fc2bb6473a09"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceModifier(action_config, g, context, target, modifiers.Stunned, 30))
			}

			return transactions.Build()
		},
	})
}
