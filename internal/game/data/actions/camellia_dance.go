package actions

import (
	"fmt"
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CamelliaDance = MakeCamelliaDance()

func MakeCamelliaDance() game.Action {
	config := game.ActionConfig{
		Name:        "Camellia Dance",
		Description: "Hits 3 times. Bypasses Protect. Always crits.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(25),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsTai),
		Cost:        game.Ptr(0),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(4)),
		CritMod:     1.5,
	}
	return game.Action{
		ID:              uuid.MustParse("c2ff8167-941a-4c2b-844f-e3f5bb7d738b"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				conf, _ := game.GetActiveActionConfig(g, config)
				crit_result := game.MakeCriticalCheck(conf)
				fmt.Printf("%+v\n", crit_result)
				damage_config := game.NewDamageConfig(crit_result.Ratio, game.RandomDamageFactor())
				damage_config.Repeat = true
				damage_config.RepeatMax = 3
				damage_config.IgnoreProtect = true
				damages := game.NewDamage(conf, damage_config)
				transactions = append(
					transactions,
					game.MakeDamageTransactions(context, damages)...,
				)

				return transactions
			},
		},
	}
}
