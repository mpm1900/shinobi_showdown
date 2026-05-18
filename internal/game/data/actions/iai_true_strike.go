package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var IaiTrueStrike = MakeIaiTrueStrike()

func MakeIaiTrueStrike() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Iai: True Strike",
		Description: "Never misses.",
		Nature:      game.Ptr(game.NsTai),
		Power:       game.Ptr(85),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(60),
		Jutsu:       game.Bukijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("d5dd8925-ae63-47e9-a2d0-f0e448129ca7"),
		Config: config,
	})
}
