package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var OneThousandNeedles = MakeOneThousandNeedles()

func MakeOneThousandNeedles() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "1000 Needles",
		Description: "Hits 2-5 times. High critical chance.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(25),
		Stat:        game.Ptr(game.StatChakraAttack),
		Nature:      game.Ptr(game.NsIce),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
	})

	config.CritStage = game.Ptr(1)

	return game.Action{
		ID:              uuid.MustParse("58c829b9-aa81-4a44-84c7-73cf08501e48"),
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
