package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var PunishingFire = MakePunishingFire()

func MakePunishingFire() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Punishing Fire",
		Description: "Double power if target is statused.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("91b571af-0f7b-42ef-8275-fec11e52c372"),
		Config: config,
		MapConfig: func(g game.Game, ctx game.Context, config game.ActionConfig) game.ActionConfig {
			targets := g.GetTargets(ctx)
			for _, target := range targets {
				if target.Statused {
					config.Power = game.Ptr(*config.Power * 2)
					break
				}
			}
			return config
		},
	})
}
