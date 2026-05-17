package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WhirlwindKick = MakeWhirlwindKick()

func MakeWhirlwindKick() game.Action {
	ID := uuid.MustParse("b23ace96-eb09-5bf7-b884-7ef8e8fc544d")
	config := game.ActionConfig{
		Name:        "Whirlwind Kick",
		Description: "High critical chance.",
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Nature:      game.Ptr(game.NsTai),
		Cost:        game.Ptr(0),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(2)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
