package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var BlackNeedle = MakeBlackNeedle()

func MakeBlackNeedle() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Black Needle",
		Description: "Lowers target's Chakra Attack.",
		Nature:      game.Ptr(game.NsYinYang),
		Accuracy:    game.Ptr(95),
		Power:       game.Ptr(75),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("74cecc3b-3297-4f79-be1a-6be167e34ac0"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				ctx := game.MakeContextForActor(target)
				ctx.SourceActorID = context.SourceActorID
				mod := modifiers.ChakraAttackDownTarget
				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, ctx)
				transactions = append(transactions, transaction)
			}

			return transactions
		},
	})
}
