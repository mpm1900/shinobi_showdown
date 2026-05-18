package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var Chidori = MakeChidori()

func MakeChidori() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Chidori",
		Description: "Deals 20% recoil damage.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(95),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(50),
		Recoil:      game.Ptr(0.2),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("c1502330-764c-56f8-9c9e-f41b933a90f0"),
		Config: config,
	})
}
