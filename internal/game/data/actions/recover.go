package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Recover = MakeRecover()

func MakeRecover() game.Action {
	return makeNoneAction(
		uuid.MustParse("c0756ddc-2611-5eef-82cc-c2bc03f9f01c"),
		makeStatusConfig(game.ActionConfig{
			Name:        "Recover",
			Description: "Heals the user for up to half of their HP.",
			Nature:      game.Ptr(game.NsYang),
			Cost:        game.Ptr(30),
			Jutsu:       game.Senjutsu,
		}),
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			action_config, _ := game.GetActiveActionConfig(g)
			ctx := game.MakeContextForActor(source)
			heal := game.NewHeal(action_config, 0.5)
			transactions.Push(game.MakeTransaction(heal, ctx))

			return transactions.Build()
		},
	)
}
