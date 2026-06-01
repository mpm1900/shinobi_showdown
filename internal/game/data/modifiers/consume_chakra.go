package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var consumeChakraID = uuid.MustParse("01bd3c19-92d2-5424-8eb0-c6132fe31062")

var ConsumeChakraTrigger game.Trigger = game.Trigger{
	ID:         uuid.MustParse("88e0a742-50c1-53b3-a873-4c3945fdd270"),
	ModifierID: consumeChakraID,
	On:         game.OnKill,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.SourceIsAlive,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			mutation := mutations.AddModifiers(false, ChakraAttackUpSource)
			transaction := game.MakeTransaction(mutation, context)
			transactions = append(transactions, transaction)

			return transactions
		},
	},
}

var ConsumeChakra game.Modifier = game.Modifier{
	ID:          consumeChakraID,
	GroupID:     &consumeChakraID,
	Icon:        "consume_chakra",
	Name:        "Consume Chakra",
	Description: "On kill: gain Chakra Attack Up.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&consumeChakraID),
	},
	Triggers: []game.Trigger{
		ConsumeChakraTrigger,
	},
}
