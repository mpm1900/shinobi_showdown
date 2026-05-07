package actions

import (
	"fmt"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var KamuiCounter = MakeKamuiCounter()

func MakeKamuiCounter() game.Action {
	ID := uuid.MustParse("09b013e2-436d-4bd4-bf0a-9b2827b9c131")
	config := game.ActionConfig{
		Name:        "Kamui: Counter",
		Description: "Fails unless the target has a pending attack. +1 priority.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsYin),
		Cost:        game.Ptr(0),
		TargetCount: game.Ptr(1),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}
	action := makeAttack(AttackConfig{
		ID:       ID,
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP1),
	})
	delta := action.Delta
	action.Delta = func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
		transactions := []game.GameTransaction{}

		// loop through every taerget and every action to see if they have a pending attack, if not, this attack fails
		for _, target := range g.GetTargets(context) {
			var found_action *game.Transaction[game.Action] = nil
			for _, action := range g.Actions {
				if action.Context.SourceActorID != nil && *action.Context.SourceActorID == target.ID {
					// attacking moves only
					if action.Mutation.Config.Power != nil {
						found_action = &action
						break
					}

				}
			}

			if found_action == nil {
				// then target has no action
				log := game.NewLog(fmt.Sprintf("%s failed.", config.Name))
				log_mut := game.AddLogs(log)
				log_tx := game.MakeTransaction(log_mut, game.NewContext())
				transactions = append(transactions, log_tx)
				return transactions
			}
		}

		return delta(p, g, context)
	}

	return action
}
