package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var KamuiSlash = MakeKamuiSlash()

func MakeKamuiSlash() game.Action {
	ID := uuid.MustParse("76fce93b-68d0-43d3-9ffb-c4560d652950")

	config := game.ActionConfig{
		Name:        "Kamui: Slash",
		Description: "Never misses.",
		Nature:      game.Ptr(game.NsYin),
		Power:       game.Ptr(85),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(60),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
	})
}
