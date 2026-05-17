package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var IronBody = MakeIronBody()

func MakeIronBody() game.Action {
	config := game.ActionConfig{
		Name:        "Iron Body",
		Nature:      game.Ptr(game.NsEarth),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Attack and Defense and lowers Speed.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("36a1a65f-cf89-44e6-9fa5-264636bf9066"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.DefenseUpSource, modifiers.SpeedDownSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
