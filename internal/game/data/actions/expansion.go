package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Expansion = MakeExpansion()

func MakeExpansion() game.Action {
	return makeNoneStatus(
		uuid.MustParse("94c7641b-c089-4c38-ae4d-56869f3d9ca6"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Expansion",
			Nature:      game.Ptr(game.NsYang),
			Jutsu:       game.Taijutsu,
			Description: "Raises the user's Attack and Chakra Attack stats.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.ChakraAttackUpSource)
			transaction := game.MakeTransaction(mutation, context)
			transactions.PushOne(transaction)

			return transactions.Build()
		},
	)
}
