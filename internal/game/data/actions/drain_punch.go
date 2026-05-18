package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var DrainPunch = MakeDrainPunch()

func MakeDrainPunch() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Drain Punch",
		Description: "User recovers 50% of damage dealt.",
		Nature:      game.Ptr(game.NsTai),
		Accuracy:    game.Ptr(75),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Recoil:      game.Ptr(-0.5),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("caa1cb26-54b3-4c9d-a305-0354e51b056e"),
		Config: config,
	})
}
