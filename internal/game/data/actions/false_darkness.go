package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var FalseDarkness = MakeFalseDarkness()

func MakeFalseDarkness() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("99338b50-de10-4747-9e41-847677db4ca0"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "False Darkness",
			Description: "Grants the user Lightning nature until end of turn.",
			Nature:      game.Ptr(game.NsLightning),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(95),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(30),
			Jutsu:       game.Ninjutsu,
		}),
		BeforeAttack: func(g game.Game, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			add_mut := mutations.AddModifiers(false, modifiers.AddNature(game.NsLightning, 0))
			add_tx := game.MakeTransaction(add_mut, context)
			transactions.Push(add_tx)

			return transactions.Build()
		},
	})
}
