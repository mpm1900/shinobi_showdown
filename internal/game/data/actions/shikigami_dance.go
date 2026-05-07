package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ShikigamiDance = MakeShikigamiDance()

func MakeShikigamiDance() game.Action {
	config := game.ActionConfig{
		Name:        "Shikigami Dance",
		Nature:      game.Ptr(game.NsWind),
		Jutsu:       game.Taijutsu,
		Description: "Raises the user's Chakra Defense by 2 stages.",
	}
	return game.Action{
		ID:              uuid.MustParse("ba9be5cb-607a-4f91-8830-7f00eaf4ea16"),
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mod := modifiers.ChakraDefenseUp2Source
				mutation := mutations.AddModifiers(false, mod)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
