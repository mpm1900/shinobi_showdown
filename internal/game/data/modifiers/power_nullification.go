package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var powerNullificationID = uuid.MustParse("eb04a402-af2a-46f8-9490-bed544b46ef3")

var PowerNullification = game.Modifier{
	ID:          powerNullificationID,
	Name:        "Power Nullification",
	Icon:        "power_nullification",
	Description: "While attacking or being attacked: Reset all stat stages.",
	Show:        true,
	GroupID:     &powerNullificationID,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		{
			ModifierGroupID: &powerNullificationID,
			Mutation: game.Mutation[game.Game, game.Actor, game.Actor]{
				Priority: game.MutPriorityStagedZero,
				Filter:   game.ActiveTransactionFilter,
				Delta: func(p game.Game, a game.Actor, c game.Context) game.Actor {
					a.Stages[game.StatAccuracy] = 0
					a.Stages[game.StatAttack] = 0
					a.Stages[game.StatChakraAttack] = 0
					a.Stages[game.StatChakraDefense] = 0
					a.Stages[game.StatDefense] = 0
					a.Stages[game.StatEvasion] = 0
					a.Stages[game.StatHP] = 0
					a.Stages[game.StatSpeed] = 0
					a.Stages[game.StatStamina] = 0

					return a
				},
			},
		},
	},
}
