package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var FlameBullet = MakeFlameBullet()

func MakeFlameBullet() game.Action {
	ID := uuid.MustParse("46389f19-31f5-4814-b8ab-32a22be9258f")

	config := game.ActionConfig{
		Name:        "Flame Bullet",
		Description: "+1 priority.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:       ID,
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
