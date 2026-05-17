package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ChannelChakra = MakeChannelChakra()

func MakeChannelChakra() game.Action {
	config := game.ActionConfig{
		Name:        "Channel Chakra",
		Nature:      game.Ptr(game.NsYang),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Chakra Attack and Chakra Defense stats.",
		TargetType:  game.TargetActorID,
	}
	return game.Action{
		ID:              uuid.MustParse("403a4bc9-d2fe-4604-a549-6f5e6c7f8dc8"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}
				mutation := mutations.AddModifiers(false, modifiers.ChakraAttackUpSource, modifiers.ChakraDefenseUpSource)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
