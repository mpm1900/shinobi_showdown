package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var actionBoostID = uuid.MustParse("ff9387f1-ee1a-4162-b5a6-ce2b47498bc5")
var ActionBoost = MakeActionBoost()

func MakeActionBoost() game.Action {
	config := game.ActionConfig{
		Name:        "Action Boost",
		Nature:      game.Ptr(game.NsTai),
		Jutsu:       game.Ninjutsu,
		Description: "Target acts after this action.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              actionBoostID,
		Config:          config,
		TargetPredicate: game.ComposeAF(game.TeamFilter, game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.TargetLengthFilter(1),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityP5,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					mutation := mutations.BoostActionPriority(target)
					transaction := game.MakeTransaction(mutation, context)
					transactions = append(transactions, transaction)
				}

				return transactions
			},
		},
	}
}

// proxies
var FlyingRaijinAssist = MakeFlyingRaijinAssist()

func MakeFlyingRaijinAssist() game.Action {
	action := MakeActionBoost()
	action.ID = uuid.MustParse("57bcb90e-df7b-4bfd-92f2-2d40847afedb")
	action.Config.Name = "Flying Raijin: Assist"
	return action
}
