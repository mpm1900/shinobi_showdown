package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var GiantRasengan = MakeGiantRasengan()
var RasenganRecharge = MakeRasenganRecharge()

func MakeGiantRasengan() game.Action {
	ID := uuid.MustParse("e0874a45-2f62-5544-a4a2-f440644407db")
	config := game.ActionConfig{
		Name:        "Giant Rasengan",
		Description: "Powerful chakra attack. Must recharge the next turn.",
		Nature:      game.Ptr(game.NsSage),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(150),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(100),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}
	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, context, _ game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			recharge := mutations.QueueAction(RasenganRecharge.ID, context)
			transactions = append(transactions, game.MakeTransaction(recharge, context))

			return transactions
		},
	})
}

func MakeRasenganRecharge() game.Action {
	config := game.ActionConfig{
		Name:       "Recharging...",
		LogSuccess: game.Ptr("$source$ must recharge."),
	}
	return game.Action{
		ID:              uuid.MustParse("2eaa6398-06a5-56fe-b90d-e9db6f044744"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.PositionsLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				return transactions
			},
		},
	}
}
