package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var Command = MakeCommand()

func MakeCommand() game.Action {
	config := game.ActionConfig{
		Name:        "Command",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Taijutsu,
		Description: "Ally uses their last used action again.",
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
	}
	return game.Action{
		ID:              uuid.MustParse("e66cbb4f-ec40-4e40-b88d-ff8a47fecfb4"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				targets := g.GetTargets(context)
				for _, target := range targets {
					if target.LastUsedActionTX == nil {
						continue
					}

					action_mut := mutations.PushExtraAction(target.LastUsedActionTX.Mutation, target.LastUsedActionTX.Context)
					action_tx := game.MakeTransaction(action_mut, target.LastUsedActionTX.Context)
					transactions = append(transactions, action_tx)
				}

				return transactions
			},
		},
	}
}
