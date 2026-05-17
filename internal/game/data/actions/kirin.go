package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Kirin = MakeKirin()

func MakeKirin() game.Action {
	ID := uuid.MustParse("d55c8221-fc03-4ae0-9737-cb5c7db88f73")

	config := game.ActionConfig{
		Name:        "Kirin",
		Description: "30% chance to paralyze the target. Always crits.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(70),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(50),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(4)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceParalysis(config, g, context, target, 30)...)
			}

			return transactions
		},
	})
}
