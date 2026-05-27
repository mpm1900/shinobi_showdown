package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var LuckyStrikes = MakeLuckyStrikes()

func MakeLuckyStrikes() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:     "Lucky Strikes",
		Accuracy: game.Ptr(80),
		Power:    game.Ptr(10),
		Stat:     game.Ptr(game.StatAttack),
		Nature:   game.Ptr(game.NsTai),
		Cost:     game.Ptr(30),
		Jutsu:    game.Taijutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("4ac4894c-2ff3-5142-b087-a8924837cefc"),
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
				damage_config := game.NewDamageConfig(game.RandomDamageFactor())
				damage_config.Repeat = true
				damage_config.RepeatMax = 999
				damages := game.DamageCoreMutation(conf, damage_config)
				transactions = append(
					transactions,
					game.MakeDamageTransactions(context, damages)...,
				)

				return transactions
			},
		},
	}
}
