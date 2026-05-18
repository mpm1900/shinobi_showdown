package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var PaperBomb = MakePaperBomb()

func MakePaperBomb() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Paper Bomb",
		Description: "Damage is based on the user's Chakra Defense rather than Chakra Attack.",
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(70),
		Stat:        game.Ptr(game.StatChakraDefense),
		Nature:      game.Ptr(game.NsExplosion),
		Cost:        game.Ptr(0),
		Jutsu:       game.Fuinjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("414f82a3-63e8-45b5-a398-7c8d15519552"),
		Config: config,
	})
}
