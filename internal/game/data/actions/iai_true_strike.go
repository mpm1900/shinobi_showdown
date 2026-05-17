package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var IaiTrueStrike = MakeIaiTrueStrike()

func MakeIaiTrueStrike() game.Action {
	ID := uuid.MustParse("d5dd8925-ae63-47e9-a2d0-f0e448129ca7")

	config := game.ActionConfig{
		Name:        "Iai: True Strike",
		Description: "Never misses.",
		Nature:      game.Ptr(game.NsTai),
		Power:       game.Ptr(85),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(60),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Bukijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
