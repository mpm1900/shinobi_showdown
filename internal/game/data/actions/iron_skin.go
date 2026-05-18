package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var IronSkin = MakeIronSkin()

func MakeIronSkin() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Iron Skin",
		Nature:      game.Ptr(game.NsEarth),
		Jutsu:       game.Ninjutsu,
		Description: "Raises the user's Defense stats by 2 stages.",
	})

	return game.Action{
		ID:              uuid.MustParse("4f70f329-a1f4-4e09-aa36-9bd4bc47198c"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, modifiers.DefenseUp2Source)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
