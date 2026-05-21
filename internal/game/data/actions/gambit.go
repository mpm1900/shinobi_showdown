package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Gambit = MakeGambit()

func MakeGambit() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Gambit",
		Description: "Attacks before switching. Fails unless target is switching.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Taijutsu,
	})

	action := makeAttack(AttackConfig{
		ID:       uuid.MustParse("68532c7e-85a5-49fd-bfa2-7d9b1e927c08"),
		Config:   config,
		Priority: game.Ptr(game.ActionPrioritySwitch),
	})

	action.Config.SubPriority = 1
	delta := action.Delta
	action.Delta = func(p, g game.Game, context game.Context) []game.GameTransaction {
		targets := g.GetTargets(context)
		is_switch := true
		for _, target := range targets {
			_, ok := g.FindQueuedAction(func(tx game.Transaction[game.Action]) bool {
				if tx.Context.SourceActorID == nil {
					is_switch = false
				}
				return *tx.Context.SourceActorID == target.ID && tx.Mutation.Config.Switch
			})
			if !ok {
				is_switch = false
			}
		}

		if !is_switch {
			log := game.NewLogContext("$action$ failed.", context)
			log_tx := game.MakeTransaction(game.AddLogs(log), context)
			return []game.GameTransaction{
				log_tx,
			}
		}
		return delta(p, g, context)
	}

	return action
}
