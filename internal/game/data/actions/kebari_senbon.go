package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var KebariSenbon = MakeKebariSenbon()

func MakeKebariSenbon() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Kebari Senbon",
		Description: "Hits 2-5 times. High critical chance.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(25),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsYang),
		Cost:        game.Ptr(50),
		Jutsu:       game.Bukijutsu,
	})

	config.CritStage = game.Ptr(1)

	return game.Action{
		ID:              uuid.MustParse("0de3affc-7513-41b0-8622-c603ccb8ee8a"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				action_config, _ := game.GetActiveActionConfig(g, config)
				damage_config := game.NewDamageConfig(game.RandomDamageFactor())
				damage_config = MakeRepeats(damage_config, 2, 5, g, context)
				damage := game.DamageCoreMutation(action_config, damage_config)
				transactions = append(
					transactions,
					game.MakeTransaction(damage, context),
				)

				return transactions
			},
		},
	}
}
