package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Amaterasu = MakeAmaterasu()

func MakeAmaterasu() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Amaterasu",
		Description: "Burns target.",
		Nature:      game.Ptr(game.NsYin),
		Stat:        game.Ptr(game.StatChakraAttack),
		Power:       game.Ptr(20),
		Accuracy:    game.Ptr(90),
		Cost:        game.Ptr(30),
		Jutsu:       game.Genjutsu,
	})

	return makeAttack(AttackConfig{
		ID:     uuid.MustParse("d103e605-9381-52fd-9cb8-450b7315a9f9"),
		Config: config,
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			for _, target := range g.GetTargets(context) {
				transactions.Push(modifiers.ApplyBurn(action_config, g, target, context))
			}

			return transactions.Build()
		},
	})
}
