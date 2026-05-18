package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var RaigoFist = MakeRaigoFist()

func MakeRaigoFist() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Raigo Fist",
		Description: "Power is increased +50 for each time user has been attacked this battle.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(50),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("24f3776b-1c64-47da-ac4b-68298fb80c05"),
		Config: config,
		MapConfig: func(g game.Game, context game.Context, config game.ActionConfig) game.ActionConfig {
			source, ok := g.GetSource(context)
			if !ok {
				return config
			}

			power := *config.Power + (source.HitCount * 50)
			config.Power = game.Ptr(power)
			return config
		},
	})
}
