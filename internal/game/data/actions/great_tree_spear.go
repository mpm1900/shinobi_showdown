package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GreatTreeSpear = MakeGreatTreeSpear()

func MakeGreatTreeSpear() game.Action {
	ID := uuid.MustParse("00888b4b-973f-5bf5-9a41-bba1c9b629b8")
	config := game.ActionConfig{
		Name:        "Great Tree Spear",
		Description: "30% chance to poison target.",
		Nature:      game.Ptr(game.NsWood),
		Accuracy:    game.Ptr(70),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(90),
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
				transactions = append(transactions, modifiers.ChancePoison(config, g, context, target, 30)...)
			}

			return transactions
		},
	})
}
