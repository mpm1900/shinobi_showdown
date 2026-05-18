package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Tailwind = MakeTailwind()

func MakeTailwind() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Tailwind",
		Nature:      game.Ptr(game.NsWind),
		Jutsu:       game.Ninjutsu,
		Description: "Doubles the speed of the user's party for 4 turns. This effect does not stack.",
	})

	return game.Action{
		ID:              uuid.MustParse("f0e7a99d-93ff-502a-b07c-6479f9a1fc30"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if checkPlayerHasModifier(g, context, modifiers.Tailwind.ID) {
					log := game.NewLogContext("$action$ failed.", context)
					log_tx := game.MakeTransaction(game.AddLogs(log), context)
					return append(transactions, log_tx)
				}

				mod := modifiers.Tailwind
				mod.Duration = 5
				mod.Icon = "tailwind"
				mutation := mutations.AddModifiers(false, mod)
				context.ParentActorID = nil
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
