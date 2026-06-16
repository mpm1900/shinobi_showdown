package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var GiantRasengan = MakeGiantRasengan()
var RasenganRecharge = MakeRasenganRecharge()

func MakeGiantRasengan() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("e0874a45-2f62-5544-a4a2-f440644407db"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Giant Rasengan",
			Description: "Powerful chakra attack. Must recharge the next turn.",
			Nature:      game.Ptr(game.NsSage),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(150),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(100),
			Jutsu:       game.Ninjutsu,
		}),
		OnSuccess: func(g game.Game, context, _ game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			recharge := mutations.QueueAction(RasenganRecharge.ID, context)
			transactions = append(transactions, game.MakeTransaction(recharge, context))

			return transactions
		},
	})
}

func MakeRasenganRecharge() game.Action {
	return makeNoneAction(
		uuid.MustParse("2eaa6398-06a5-56fe-b90d-e9db6f044744"),
		game.ActionConfig{
			Name:       "Recharging...",
			LogSuccess: game.Ptr("$source$ must recharge."),
		},
		func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			return game.NewTransactionBuilder().Build()
		},
	)
}
