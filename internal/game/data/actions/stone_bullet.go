package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var StoneBullet = MakeStoneBullet()

func MakeStoneBullet() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Stone Bullet",
		Description: "10% chance to stun target.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("37ff4e8e-14e1-4d00-bde9-3f834b97cb73"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceModifier(action_config, g, context, target, modifiers.Stunned, 10))
			}

			return transactions.Build()
		},
	})
}
