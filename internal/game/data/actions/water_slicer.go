package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WaterSlicer = MakeWaterSlicer()

func MakeWaterSlicer() game.Action {
	ID := uuid.MustParse("b6f7440c-8f3a-44a8-9412-ea6db426ccd3")
	config := game.ActionConfig{
		Name:        "Water Slicer",
		Description: "+1 priority. High critical hit chance.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(40),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(1),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(1)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:       ID,
		Config:   config,
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
