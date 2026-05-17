package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var MudWall = MakeMudWall()

func MakeMudWall() game.Action {
	config := game.ActionConfig{
		Name:        "Mud Wall",
		Nature:      game.Ptr(game.NsEarth),
		Jutsu:       game.Ninjutsu,
		Description: "User's team takes 50% less physical damage for 5 turns.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("8ddc221c-ebb2-47cf-bbe4-09da335eb70b"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if checkPlayerHasModifier(g, context, modifiers.PhysicalDamageDownTeam.ID) {
					log := game.NewLogContext("$action$ failed.", context)
					log_tx := game.MakeTransaction(game.AddLogs(log), context)
					return append(transactions, log_tx)
				}

				mod := modifiers.PhysicalDamageDownTeam
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
