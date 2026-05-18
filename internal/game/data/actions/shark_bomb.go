package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var SharkBomb = MakeSharkBomb()

func MakeSharkBomb() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Shark Bomb",
		Description: "Deals 30% recoil damage.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(70),
		Recoil:      game.Ptr(0.3),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("28f4a572-b566-40dc-ae22-2f5e22cc5229"),
		Config: config,
	})
}
