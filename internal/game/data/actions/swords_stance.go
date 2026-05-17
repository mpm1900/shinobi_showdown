package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SwordsStance = MakeSwordsStance()

func MakeSwordsStance() game.Action {
	config := game.ActionConfig{
		Name:        "Swords Stance",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Taijutsu,
		Description: "Raises the user's Physical Attack by 2 stages.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("cdda818c-edac-4de4-99e8-d0890fcc9214"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mod := modifiers.AttackUp2Source
				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
