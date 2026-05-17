package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Rest = MakeRest()

func MakeRest() game.Action {
	config := game.ActionConfig{
		Name:        "Rest",
		Description: "Heals user 100%. User falls asleep.",
		Nature:      game.Ptr(game.NsYang),
		TargetCount: game.Ptr(0),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(30),
		Jutsu:       game.Senjutsu,
	}

	return game.Action{
		ID:              uuid.MustParse("64372d78-3c0b-4d57-a71c-8bdf0e2aedc0"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				source_ctx := game.MakeContextForActor(source)
				conf, _ := game.GetActiveActionConfig(g, config)
				heal := game.NewHeal(conf, 1)
				sleep := modifiers.ApplySleep(config, g, source, context)
				transactions = append(transactions, sleep...)
				transactions = append(
					transactions,
					game.MakeTransaction(heal, source_ctx),
				)

				return transactions
			},
		},
	}
}
