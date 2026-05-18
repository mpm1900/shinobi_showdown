package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var FlyingSwallow = MakeFlyingSwallow()

func MakeFlyingSwallow() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Flying Swallow",
		Nature:      game.Ptr(game.NsWind),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Attack and Chakra Attack stats.",
	})

	return game.Action{
		ID:              uuid.MustParse("497c3176-dc06-4762-ba5f-cc029c54f258"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.AttackUpSource, modifiers.ChakraAttackUpSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
