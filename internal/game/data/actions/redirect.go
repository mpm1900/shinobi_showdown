package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Redirect = MakeRedirect()

func MakeRedirect() game.Action {
	action := makeNoneAction(
		uuid.MustParse("3d0b6e04-f5f0-50db-9eb6-4aede4c11701"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Redirect",
			Nature:      game.Ptr(game.NsYin),
			Jutsu:       game.Genjutsu,
			Description: "Changes the target of single-target enemy actions.",
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			mutation := mutations.RedirectSingleTargetEnemyActions(source)
			transaction := game.MakeTransaction(mutation, context)
			transactions.Push(transaction)

			return transactions.Build()
		},
	)

	action.Priority = game.ActionPriorityP2
	return action
}
