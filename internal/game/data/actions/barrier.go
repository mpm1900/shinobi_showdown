package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Barrier = MakeBarrier()

func MakeBarrier() game.Action {
	config := game.ActionConfig{
		Name:        "Barrier",
		Nature:      game.Ptr(game.NsYin),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		Description: "Protects the user's team from multi-target actions. +3 priority, 1 turn cooldown.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("fd909d15-b6c4-4892-a5d2-9af752d52cc8"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP3,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
				game.SourceIsActionOffCooldown,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.SpreadProtected)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
