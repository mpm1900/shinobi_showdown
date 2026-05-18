package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var KamuiSlash = MakeKamuiSlash()

func MakeKamuiSlash() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Kamui: Slash",
		Description: "Never misses.",
		Nature:      game.Ptr(game.NsYin),
		Power:       game.Ptr(85),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(60),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("76fce93b-68d0-43d3-9ffb-c4560d652950"),
		Config: config,
	})
}
