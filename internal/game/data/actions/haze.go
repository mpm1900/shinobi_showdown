package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Haze = MakeHaze()

func MakeHaze() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Haze",
		Nature:      game.Ptr(game.NsWater),
		Jutsu:       game.Ninjutsu,
		Description: "Nullifies all stat stage changes for 5 turns.",
	})

	return game.Action{
		ID:              uuid.MustParse("63db7718-b73b-5f31-8b1f-c2dfa5bd5c65"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				for _, tx := range g.Modifiers {
					if tx.Context.SourcePlayerID == nil {
						continue
					}

					if *tx.Context.SourcePlayerID == *context.SourcePlayerID && tx.Mutation.ID == modifiers.Haze.ID {
						log := game.MakeGameLog("$action$ failed.", context, 0)
						log_tx := game.MakeTransaction(game.AddLogs(log), context)
						return append(transactions, log_tx)
					}

				}

				mod := modifiers.Haze
				mod.Duration = 5
				mutation := mutations.AddModifiers(false, mod)
				context.ParentActorID = nil
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
