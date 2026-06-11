package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Kirin = MakeKirin()

func MakeKirin() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Kirin",
		Description: "30% chance to paralyze the target. Always crits.",
		Nature:      game.Ptr(game.NsLightning),
		Accuracy:    game.Ptr(70),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(50),
		Jutsu:       game.Ninjutsu,
		CritStage:   game.Ptr(4),
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("d55c8221-fc03-4ae0-9737-cb5c7db88f73"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Concat(modifiers.ChanceParalysis(action_config, g, context, target, 30))
			}

			return transactions.Build()
		},
	})
}
