package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var WaterWall = MakeWaterWall()

func MakeWaterWall() game.Action {
	return makeNoneAction(
		uuid.MustParse("e615a56a-2d0b-46ec-99cc-66553c5ed6c1"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Water Wall",
			Nature:      game.Ptr(game.NsWater),
			Jutsu:       game.Ninjutsu,
			Description: "User's team takes 50% less chakra damage for 5 turns.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mod := modifiers.ChakraDamageDownTeam
			mod.Duration = 5

			if checkPlayerHasModifier(g, context, mod.ID) {
				log := game.NewLogContext("$action$ failed.", context)
				log_tx := game.MakeTransaction(game.AddLogs(log), context)
				transactions.Push(log_tx)
				return transactions.Build()
			}

			mutation := mutations.AddModifiers(false, mod)
			context.ParentActorID = nil
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)
			return transactions.Build()
		},
	)
}
