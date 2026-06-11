package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var RockFist = MakeRockFist()

func MakeRockFist() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Rock Fist",
		Description: "30% chance to paralyze the target.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("6e015595-18df-4f2b-b4da-cf97863a3f4e"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceParalysis(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
