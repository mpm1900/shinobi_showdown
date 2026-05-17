package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var WaterWall = MakeWaterWall()

func MakeWaterWall() game.Action {
	config := game.ActionConfig{
		Name:        "Water Wall",
		Nature:      game.Ptr(game.NsWater),
		Jutsu:       game.Ninjutsu,
		Description: "User's team takes 50% less chakra damage for 5 turns.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("e615a56a-2d0b-46ec-99cc-66553c5ed6c1"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				if checkPlayerHasModifier(g, context, modifiers.ChakraDamageDownTeam.ID) {
					log := game.NewLogContext("$action$ failed.", context)
					log_tx := game.MakeTransaction(game.AddLogs(log), context)
					return append(transactions, log_tx)
				}

				mod := modifiers.ChakraDamageDownTeam
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
