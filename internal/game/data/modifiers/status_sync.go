package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var statusSyncID = uuid.MustParse("286eb78f-880a-4514-8e05-fa237292cb32")
var StatusSyncTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: statusSyncID,
	On:         game.OnStatusAdd,
	Check:      game.Match__TargetActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			targets := g.GetTargets(context)
			for _, target := range targets {
				if !target.Statused {
					continue
				}
				if target.ID == source.ID {
					break
				}
				if target.Burned {
					transactions.Push(ApplyBurn(game.ActionConfig{}, g, source, game.NewContext()))
				}
				if target.Paralyzed {
					transactions.Push(ApplyParalysis(game.ActionConfig{}, g, source, game.NewContext()))
				}
				if target.Poisoned {
					transactions.Push(ApplyPoison(game.ActionConfig{}, g, source, game.NewContext()))
				}
				if target.Sleeping {
					transactions.Push(ApplySleep(game.ActionConfig{}, g, source, game.NewContext()))
				}
				break
			}

			return transactions.Build()
		},
	},
}

var StatusSync = game.Modifier{
	ID:          statusSyncID,
	GroupID:     &statusSyncID,
	Icon:        "status_sync",
	Name:        "Status Sync",
	Description: "On status: source receives the same status.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&statusSyncID),
	},
	Triggers: []game.Trigger{
		StatusSyncTrigger,
	},
}
