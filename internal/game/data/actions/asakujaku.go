package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Asakujaku = MakeAsakujaku()

func MakeAsakujaku() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("5c7a660c-7d53-48a2-94c3-c2b683fed948"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "6th Gate: Asakujaku",
			Description: "Deals 30% recoil damage. 10% chance to burn target. High crit chance.",
			Nature:      game.Ptr(game.NsFire),
			Accuracy:    game.Ptr(95),
			Power:       game.Ptr(120),
			Stat:        game.Ptr(game.StatAttack),
			Recoil:      game.Ptr(0.3),
			Jutsu:       game.Taijutsu,
			CritStage:   game.Ptr(1),
		}),
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			for _, target := range g.GetTargets(context) {
				transactions.Concat(modifiers.ChanceBurn(action_config, g, context, target, 10))
			}

			return transactions.Build()
		},
	})
}
