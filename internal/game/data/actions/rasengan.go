package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Rasengan = MakeRasengan()

func MakeRasengan() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("054eb97a-cd6f-4428-8f54-96d9b6b33bfa"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:     "Rasengan",
			Nature:   game.Ptr(game.NsSage),
			Accuracy: game.Ptr(90),
			Power:    game.Ptr(90),
			Stat:     game.Ptr(game.StatChakraAttack),
			Cost:     game.Ptr(50),
			Jutsu:    game.Ninjutsu,
		}),
	})
}
