package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var flyingLotusID = uuid.MustParse("75d0a924-912a-4972-93fb-3c08e82bb1b3")
var FlyingLotus = MakeFlyingLotus()

func MakeFlyingLotus() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Flying Lotus",
		Description: "Must attack for 3 turns.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     flyingLotusID,
		Config: config,
		OnSuccess: func(g game.Game, context game.Context, _ game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			key := "repeats"
			repeats, ok := context.GetMeta(key)
			if !ok {
				repeats = 0
			}

			context.SetMeta(key, repeats+1)
			if repeats < 2 {
				recharge := mutations.QueueAction(flyingLotusID, context)
				transactions.PushOne(game.MakeTransaction(recharge, context))
			}

			return transactions.Build()
		},
	})
}
