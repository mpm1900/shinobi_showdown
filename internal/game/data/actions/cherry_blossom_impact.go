package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CherryBlossomImpact = MakeCherryBlossomImpact()

func MakeCherryBlossomImpact() game.Action {
	ID := uuid.MustParse("5188108e-d87b-4405-9006-ef0667ea4a62")
	config := game.ActionConfig{
		Name:        "Cherry Blossom Impact",
		Description: "30% chance to burn the target.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 30)...)
			}

			return transactions
		},
	})
}
