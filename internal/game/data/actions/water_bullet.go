package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var WaterBullet = MakeWaterBullet()

func MakeWaterBullet() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Water Bullet",
		Description: "10% chance to lower target's speed.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("31269022-d87e-4142-bae0-1235ef882112"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceModifier(action_config, g, context, target, modifiers.SpeedDownTarget, 10))
			}

			return transactions.Build()
		},
	})
}
