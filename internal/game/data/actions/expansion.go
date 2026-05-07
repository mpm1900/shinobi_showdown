package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Expansion = MakeExpansion()

func MakeExpansion() game.Action {
	config := game.ActionConfig{
		Name:        "Expansion",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Taijutsu,
		Description: "Raises the user's Attack and Defense stats.",
	}
	return game.Action{
		ID:              uuid.MustParse("94c7641b-c089-4c38-ae4d-56869f3d9ca6"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}
				mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.DefenseUpSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
