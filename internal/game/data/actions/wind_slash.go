package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var WindSlash = MakeWindSlash()

func MakeWindSlash() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Wind Slash",
		Description: "20% chance to stun target. High critical hit chance.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(60),
		Jutsu:       game.Ninjutsu,
	})
	config.CritChance = game.Ptr(getCriticalStage(1))

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("deae75a9-2943-4934-bd67-f1b773e7035f"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceModifier(config, g, context, target, modifiers.Stunned, 20)...)
			}

			return transactions
		},
	})
}
