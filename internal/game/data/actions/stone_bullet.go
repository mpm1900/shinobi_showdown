package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var StoneBullet = MakeStoneBullet()

func MakeStoneBullet() game.Action {
	ID := uuid.MustParse("37ff4e8e-14e1-4d00-bde9-3f834b97cb73")

	config := game.ActionConfig{
		Name:        "Stone Bullet",
		Description: "10% chance to stun target.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(1),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(40),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Taijutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceModifier(config, g, context, target, modifiers.Stunned, 10)...)
			}

			return transactions
		},
	})
}
