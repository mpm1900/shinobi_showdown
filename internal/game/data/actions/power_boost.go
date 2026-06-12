package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var PowerBoost = MakePowerBoost()

func MakePowerBoost() game.Action {
	action := makeSingleAction(
		uuid.MustParse("fa0a4e99-9b26-5962-9ed0-fc88a6e73cb5"),
		makeStatusConfig(game.ActionConfig{
			Name:        "Power Boost",
			Nature:      game.Ptr(game.NsSage),
			Jutsu:       game.Ninjutsu,
			Description: "Powers up target's attacks this turn.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.PowerBoosted)
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)

			return transactions.Build()
		},
	)

	action.TargetPredicate = game.ComposeAF(game.TeamFilter, game.OtherFilter, game.TargetableFilter)
	action.ActionMutation.Priority = game.ActionPriorityP5
	return action
}
