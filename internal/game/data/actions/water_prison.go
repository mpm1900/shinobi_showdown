package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var WaterPrison = MakeWaterPrison()

func MakeWaterPrison() game.Action {
	return makeSingleAction(
		uuid.MustParse("41893bef-2aad-4de6-82f6-a4f0391916d5"),
		makeStatusConfig(game.ActionConfig{
			Name:        "Water Prison",
			Nature:      game.Ptr(game.NsWater),
			Jutsu:       game.Ninjutsu,
			Description: "Traps target in a water prison that damages the user for 1/8th HP for 3 turns.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			targets := g.GetTargets(context)
			for _, target := range targets {
				context.ParentActorID = &target.ID
				mutation := mutations.AddModifiers(true, modifiers.WaterPrison)
				transaction := game.MakeTransaction(mutation, context)
				transactions.Push(transaction)
			}

			return transactions.Build()
		},
	)
}
