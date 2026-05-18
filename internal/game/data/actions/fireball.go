package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Fireball = MakeFireball()

func MakeFireball() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Fireball",
		Description: "10% chance to burn target.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("aaf5174b-f386-54b1-84c4-0c062937c770"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 10)...)
			}

			return transactions
		},
	})
}
