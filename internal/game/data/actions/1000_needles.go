package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var OneThousandNeedles = MakeOneThousandNeedles()

func MakeOneThousandNeedles() game.Action {
	config := game.ActionConfig{
		Name:        "1000 Needles",
		Description: "Hits 2-5 times. High critical chance.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(25),
		Stat:        game.Ptr(game.StatChakraAttack),
		Nature:      game.Ptr(game.NsIce),
		Cost:        game.Ptr(50),
		TargetCount: game.Ptr(1),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(1)),
		CritMod:     1.5,
	}
	return game.Action{
		ID:              uuid.MustParse("58c829b9-aa81-4a44-84c7-73cf08501e48"),
		Config:          config,
		TargetType:      game.TargetPositionID,
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
				damage_config := game.NewDamageConfig(crit_result.Ratio, game.RandomDamageFactor())
				damage_config = MakeRepeats(damage_config, 2, 5, g, context)
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
