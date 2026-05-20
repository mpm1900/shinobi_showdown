package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var FangOverFang = MakeFangOverFang()

func MakeFangOverFang() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Fang Over Fang",
		Description: "Deals 30% recoil damage.",
		Nature:      game.Ptr(game.NsYang),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(120),
		Stat:        game.Ptr(game.StatAttack),
		Recoil:      game.Ptr(0.3),
		Cost:        game.Ptr(0),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("0385152f-4956-4a2f-b289-25217e7a25bc"),
		Config: config,
	})
}
