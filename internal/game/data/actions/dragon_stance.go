package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var DragonStance = MakeDragonStance()

func MakeDragonStance() game.Action {
	return makeNoneStatus(
		uuid.MustParse("435490c1-ede2-5875-9edf-1c36d4917741"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Dragon Stance",
			Nature:      game.Ptr(game.NsTai),
			Jutsu:       game.Taijutsu,
			Description: "Raises the user's Speed and Attack stats.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.SpeedUpSource)
			transaction := game.MakeTransaction(mutation, context)
			transactions.PushOne(transaction)

			return transactions.Build()
		},
	)
}
