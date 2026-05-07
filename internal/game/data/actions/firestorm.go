package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Firestorm = MakeFirestorm()

func MakeFirestorm() game.Action {
	ID := uuid.MustParse("5756b76d-dd39-460c-b5fa-431b80200f3b")

	config := game.ActionConfig{
		Name:        "Firestorm",
		Description: "10% chance to burn target. Lowers user's Chakra attack by 2 stages.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(130),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(100),
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
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 10)...)
			}

			mod := modifiers.ChakraAttackDown2Source
			mut := mutations.AddModifiers(false, mod)
			chakraDownTx := game.MakeTransaction(mut, context)
			transactions = append(transactions, chakraDownTx)

			return transactions
		},
	})
}
