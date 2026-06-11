package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CherryBlossomImpact = MakeCherryBlossomImpact()

func MakeCherryBlossomImpact() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Cherry Blossom Impact",
		Description: "30% chance to burn the target.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("5188108e-d87b-4405-9006-ef0667ea4a62"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceBurn(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
