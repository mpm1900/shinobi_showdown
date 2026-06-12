package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var WaterSlicer = MakeWaterSlicer()

func MakeWaterSlicer() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("b6f7440c-8f3a-44a8-9412-ea6db426ccd3"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Water Slicer",
			Description: "+1 priority. High critical hit chance.",
			Nature:      game.Ptr(game.NsWater),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(40),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(0),
			Jutsu:       game.Ninjutsu,
		}),
		Priority: game.Ptr(game.ActionPriorityP1),
	})
}
