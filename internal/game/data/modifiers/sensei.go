package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var senseiID = uuid.MustParse("8bd9b3bf-0d33-4b6c-ba17-045f47b7f1d9")
var Sensei = game.Modifier{
	ID:          senseiID,
	GroupID:     &senseiID,
	Name:        "Sensei",
	Icon:        "sensei",
	Description: "Changes nature to last used attack's nature.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&senseiID),
	},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: senseiID,
			On:         game.OnActionStart,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}

					if context.ActionID == nil {
						return transactions
					}

					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					action, ok := g.FindQueuedAction(func(tx game.Transaction[game.Action]) bool {
						if tx.Context.SourceActorID == nil || context.SourceActorID == nil {
							return false
						}
						return *tx.Context.SourceActorID == *context.SourceActorID &&
							tx.Mutation.ID == *context.ActionID
					})

					if !ok {
						a, found := source.GetActionByID(g, *context.ActionID)
						if !found {
							return transactions
						}
						action = game.MakeTransaction(a, context)
					}

					if action.Mutation.Config.Nature == nil || action.Mutation.Config.Power == nil {
						return transactions
					}

					add_mut := mutations.AddModifiers(false, SetNature(*action.Mutation.Config.Nature, game.ModifierDurationInf))
					add_tx := game.MakeTransaction(add_mut, context)
					transactions = append(transactions, add_tx)

					return transactions
				},
			},
		},
	},
}
