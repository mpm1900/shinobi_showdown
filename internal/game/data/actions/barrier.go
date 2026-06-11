package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Barrier = MakeBarrier()

func MakeBarrier() game.Action {
	action := makeNoneAction(
		uuid.MustParse("fd909d15-b6c4-4892-a5d2-9af752d52cc8"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Barrier",
			Nature:      game.Ptr(game.NsYin),
			Jutsu:       game.Ninjutsu,
			Description: "Protects the user's team from multi-target actions. +3 priority.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.SpreadProtected)
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)

			return transactions.Build()
		},
	)

	action.Priority = game.ActionPriorityP3
	return action
}
