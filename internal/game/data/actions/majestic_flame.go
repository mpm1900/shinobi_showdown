package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var MajesticFlame = MakeMajesticFlame()

func MakeMajesticFlame() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Majestic Flame",
		Description: "40% chance to burn targets.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(120),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("19c9f58e-9012-417f-af58-1c09d448f0dc"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 40)...)
			}

			return transactions
		},
	})
}
