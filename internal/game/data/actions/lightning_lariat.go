package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var LightningLariat = MakeLightningLariat()

func MakeLightningLariat() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("75a444bf-cb42-4a10-8f92-6bc7de709f26"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "Lightning Lariat",
			Description: "10% chance to stun target. High crit chance.",
			Nature:      game.Ptr(game.NsLightning),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(90),
			Stat:        game.Ptr(game.StatAttack),
			Cost:        game.Ptr(50),
			Jutsu:       game.Ninjutsu,
			CritStage:   game.Ptr(1),
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			mod := modifiers.Stunned
			mod.Duration = 0

			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceModifier(action_config, g, context, target, mod, 10))
			}

			return transactions.Build()
		},
	})
}
