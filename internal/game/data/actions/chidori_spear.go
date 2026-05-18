package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var ChidoriSpear = MakeChidoriSpear()

func MakeChidoriSpear() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Chidori Spear",
		Description: "20% chance to paralyze the target.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("f89e4aba-35c8-4ddf-ab0b-bb809d5deb69"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceParalysis(config, g, context, target, 20)...)
			}

			return transactions
		},
	})
}
