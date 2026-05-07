package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var chainsOfPainID = uuid.MustParse("4735395b-d05d-4e68-9eca-2b5990de4afb")

var ChainsOfPain game.Modifier = game.Modifier{
	ID:                chainsOfPainID,
	GroupID:           &chainsOfPainID,
	Icon:              "std_chakra",
	Name:              "Chains of Pain",
	Description:       "Chakra Attack x0.75.",
	ParentDescription: "Other shinobi: Chakra Attack x0.75.",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&chainsOfPainID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.OtherFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stats[game.StatChakraAttack] = game.Round(float64(actor.Stats[game.StatChakraAttack]) * 0.75)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
