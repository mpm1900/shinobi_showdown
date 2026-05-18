package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var FlameBullet = MakeFlameBullet()

func MakeFlameBullet() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Flame Bullet",
		Description: "+1 priority.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(60),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:       uuid.MustParse("46389f19-31f5-4814-b8ab-32a22be9258f"),
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
