package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SageMode = MakeSageMode()

func MakeSageMode() game.Action {
	config := game.ActionConfig{
		Name:        "Sage Mode",
		Nature:      game.Ptr(game.NsSage),
		Jutsu:       game.Senjutsu,
		Description: "Inverts the speed of all active shinobi.",
	}
	return game.Action{
		ID:              uuid.MustParse("02796a9b-add5-5a5c-a01b-5bc6e26d0135"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPrioritySlow3,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				/*
					if checkPlayerHasModifier(g, context, modifiers.SageMode.ID) {
						log := game.NewLogContext("$action$ failed.", context)
						log_tx := game.MakeTransaction(game.AddLogs(log), context)
						return append(transactions, log_tx)
					}
				*/

				mod := modifiers.SageMode
				mod.Duration = 5
				mod.Icon = "sage_mode"
				mod.Description = "Speed is inverted."
				mutation := mutations.AddModifiers(false, mod)
				context.ParentActorID = nil
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
