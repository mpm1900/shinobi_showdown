package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Dissipate = MakeDissipate()

func MakeDissipate() game.Action {
	config := game.ActionConfig{
		Name:        "Dissipate",
		Nature:      game.Ptr(game.NsWind),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Evasion.",
	}
	return game.Action{
		ID:              uuid.MustParse("aa485c2f-920b-431a-a33d-919decdab2a4"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.EvasionUpSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
