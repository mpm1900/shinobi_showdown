package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GalePalm = MakeGalePalm()

func MakeGalePalm() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Gale Palm",
		Description: "10% chance to stun target.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(55),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(40),
		Jutsu:       game.Taijutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("0cd6cb29-5a23-41ef-94a1-348ae5c33b30"),
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
