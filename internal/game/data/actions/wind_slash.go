package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var WindSlash = MakeWindSlash()

func MakeWindSlash() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("deae75a9-2943-4934-bd67-f1b773e7035f"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Wind Slash",
			Description: "20% chance to stun target. High critical hit chance.",
			Nature:      game.Ptr(game.NsWind),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(80),
			Stat:        game.Ptr(game.StatAttack),
			Cost:        game.Ptr(60),
			Jutsu:       game.Ninjutsu,
			CritStage:   game.Ptr(1),
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceModifier(action_config, g, context, target, modifiers.Stunned, 20))
			}

			return transactions.Build()
		},
	})
}
