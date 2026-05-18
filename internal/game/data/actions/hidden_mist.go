package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var HiddenMist = MakeHiddenMist()

func MakeHiddenMist() game.Action {
	config := game.ActionConfig{
		Name:        "Hidden Mist",
		Nature:      game.Ptr(game.NsWater),
		Jutsu:       game.Ninjutsu,
		Description: "All non-water shinobi have Accuracy down x2 for 5 turns.",
	}

	return game.Action{
		ID:              uuid.MustParse("a8c7fab6-c3e0-4933-ab1a-3d05376c05b1"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if checkPlayerHasModifier(g, context, modifiers.HiddenMist.ID) {
					log := game.NewLogContext("$action$ failed.", context)
					log_tx := game.MakeTransaction(game.AddLogs(log), context)
					return append(transactions, log_tx)
				}

				mod := modifiers.HiddenMist
				mod.Show = true
				mod.Duration = 5
				filter := mod.ActorMutations[0].Filter
				mod.ActorMutations[0].Filter = func(g game.Game, a game.Actor, context game.Context) bool {
					_, ok := a.Natures[game.NsWater]
					if ok {
						return false
					}
					return filter(g, a, context)
				}
				mutation := mutations.AddModifiers(false, mod)
				context.ParentActorID = nil
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
