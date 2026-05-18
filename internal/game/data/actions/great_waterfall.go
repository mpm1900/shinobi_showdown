package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GreatWaterfall = MakeGreatWaterfall()

func MakeGreatWaterfall() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Great Waterfall",
		Description: "20% chance to stun target.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(60),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("65a8447d-4262-454d-a4ad-062993b1f8ad"),
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
