package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var LightningLigerBomb = MakeLightningLigerBomb()

func MakeLightningLigerBomb() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Lightning Liger Bomb",
		Description: "Deals double damage in Electrified Terrain.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(65),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("4b45490f-b1b9-4c76-b15f-0e8d1e5019cd"),
		Config: config,
		MapConfig: func(g game.Game, context game.Context, action_config game.ActionConfig) game.ActionConfig {
			ratio := 1.0
			state, _ := g.GetState(context)
			if state.Terrain == game.GameTerrainElectrified {
				ratio = 2.0
			}

			power := game.Round(float64(*action_config.Power) * ratio)
			action_config.Power = &power
			return action_config
		},
	})
}
