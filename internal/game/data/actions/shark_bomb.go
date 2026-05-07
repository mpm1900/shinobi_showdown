package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var SharkBomb = MakeSharkBomb()

func MakeSharkBomb() game.Action {
	ID := uuid.MustParse("28f4a572-b566-40dc-ae22-2f5e22cc5229")
	config := game.ActionConfig{
		Name:        "Shark Bomb",
		Description: "Deals 30% recoil damage.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(70),
		TargetCount: game.Ptr(1),
		Recoil:      game.Ptr(0.3),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}
	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
