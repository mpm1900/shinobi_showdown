package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SwordsStance = MakeSwordsStance()

func MakeSwordsStance() game.Action {
	return makeNoneStatus(
		uuid.MustParse("cdda818c-edac-4de4-99e8-d0890fcc9214"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Swords Stance",
			Nature:      game.Ptr(game.NsTai),
			Jutsu:       game.Taijutsu,
			Description: "Raises the user's Physical Attack by 2 stages.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.AttackUp2Source)
			transaction := game.MakeTransaction(mutation, context)
			transactions.PushOne(transaction)

			return transactions.Build()
		},
	)
}
