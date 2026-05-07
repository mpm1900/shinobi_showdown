package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Caltrops = MakeCaltrops()

func MakeCaltrops() game.Action {
	config := game.ActionConfig{
		Name:        "Caltrops",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Bukijutsu,
		Description: "Adds entry hazard for enemies.",
	}
	return game.Action{
		ID:              uuid.MustParse("e26b23d1-4f5b-4246-97a9-d29eb57b049b"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if context.SourcePlayerID == nil {
					return transactions
				}

				enemy, ok := g.GetPlayer(func(p game.Player) bool {
					return p.ID != *context.SourcePlayerID
				})
				if !ok {
					return transactions
				}

				mod_ctx := game.NewContext().WithPlayer(enemy.ID)
				mod_ctx.ParentActorID = nil

				if checkPlayerHasModifier(g, mod_ctx, modifiers.EntryHazard.ID) {
					log := game.NewLogContext("$action$ failed.", context)
					log_tx := game.MakeTransaction(game.AddLogs(log), context)
					return append(transactions, log_tx)
				}

				mod := modifiers.EntryHazard
				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, mod_ctx)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
