package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var DisarmingStrike = MakeDisarmingStrike()

func MakeDisarmingStrike() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("bfeccb7c-4d50-4575-8a95-f3357d6f81ae"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Disarming Strike",
			Description: "Target loses their held item.",
			Nature:      game.Ptr(game.NsTai),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(50),
			Stat:        game.Ptr(game.StatAttack),
			Cost:        game.Ptr(50),
			Jutsu:       game.Taijutsu,
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			targets := g.GetTargets(context)
			for _, target := range targets {
				ctx := game.MakeContextForActor(target)
				tx := game.MakeTransaction(mutations.RemoveItem, ctx)
				transactions.Push(tx)
			}

			return transactions.Build()
		},
	})
}
