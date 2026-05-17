package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WillOfTheFallen = MakeWillOfTheFallen()

func MakeWillOfTheFallen() game.Action {
	ID := uuid.MustParse("abb5de14-3093-40c6-82a6-f2bdb80cc074")
	config := game.ActionConfig{
		Name:        "Will Of The Fallen",
		Description: "Power is increased +50 for each fallen ally.",
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
			dead_ally_count := 0
			for _, ally := range g.GetActorsFilters(context, game.TeamFilter) {
				if !ally.Alive {
					dead_ally_count++
				}
			}

			power := *config.Power + (dead_ally_count * 50)
			config.Power = game.Ptr(power)
			return config
		},
	})
}
