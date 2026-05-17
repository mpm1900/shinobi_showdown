package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var PaperBomb = MakePaperBomb()

func MakePaperBomb() game.Action {
	ID := uuid.MustParse("414f82a3-63e8-45b5-a398-7c8d15519552")
	config := game.ActionConfig{
		Name:        "Paper Bomb",
		Description: "Damage is based on the user's Chakra Defense rather than Chakra Attack.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatChakraDefense),
		Nature:      game.Ptr(game.NsExplosion),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Fuinjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
