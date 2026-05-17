package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var RaigoFist = MakeRaigoFist()

func MakeRaigoFist() game.Action {
	ID := uuid.MustParse("24f3776b-1c64-47da-ac4b-68298fb80c05")
	config := game.ActionConfig{
		Name:        "Raigo Fist",
		Description: "Power is increased +50 for each time user has been attacked this battle.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(50),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(50),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
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
