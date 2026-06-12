package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Rest = MakeRest()

func MakeRest() game.Action {
	return makeNoneAction(
		uuid.MustParse("64372d78-3c0b-4d57-a71c-8bdf0e2aedc0"),
		makeNoTargetStatusConfig(game.ActionConfig{
			Name:        "Rest",
			Description: "Heals user 100%. User falls asleep.",
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

			source_ctx := game.MakeContextForActor(source)
			action_config, _ := game.GetActiveActionConfig(g)
			heal := game.NewHeal(action_config, 1)
			sleep := modifiers.ApplySleep(action_config, g, source, context)
			transactions.Concat(sleep)
			transactions.Push(game.MakeTransaction(heal, source_ctx))

			return transactions.Build()
		},
	)
}
