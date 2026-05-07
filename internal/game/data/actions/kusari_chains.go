package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var KusariChains = MakeKusariChains()

func MakeKusariChains() game.Action {
	ID := uuid.MustParse("fba5ad6e-d2ee-4b3a-b524-fc2bb6473a09")

	config := game.ActionConfig{
		Name:        "Kusari Chains",
		Description: "30% chance to stun target.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Bukijutsu,
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
				transactions = append(transactions, modifiers.ChanceModifier(config, g, context, target, modifiers.Stunned, 30)...)
			}

			return transactions
		},
	})
}
