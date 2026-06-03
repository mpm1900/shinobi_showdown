package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CamelliaDance = MakeCamelliaDance()

func MakeCamelliaDance() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Camellia Dance",
		Description: "Hits 3 times. Bypasses Protect. Always crits.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(25),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsTai),
		Cost:        game.Ptr(0),
		Jutsu:       game.Taijutsu,
	})

	config.CritStage = game.Ptr(4)

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
				transactions := game.NewTransactionBuilder()

				action_config, _ := game.GetActiveActionConfig(g, config)
				damage_config := game.NewDamageConfig(game.RandomDamageFactor())
				damage_config.Repeat = true
				damage_config.RepeatMax = 3
				damage_config.IgnoreProtect = true
				transactions.Push(game.ResolveDamageCore(action_config, damage_config, g, context))

				return transactions.Build()
			},
		},
	}
}
