package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var PhoenixFlower = MakePhoenixFlower()

func MakePhoenixFlower() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Phoenix Flower",
		Description: "Hits up-to 6 times. High critical chance.",
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(20),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsFire),
		Cost:        game.Ptr(50),
		Jutsu:       game.Bukijutsu,
	})
	config.CritStage = game.Ptr(1)

	return game.Action{
		ID:              uuid.MustParse("c6a59042-5fa2-4ec6-b83f-b705d5cd5c9e"),
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
				damage_config.RepeatMax = 6
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
