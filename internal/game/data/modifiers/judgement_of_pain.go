package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var judgementOfPainID = uuid.MustParse("8cd16d37-5336-44ae-80d0-24317038f614")

var JudgementOfPain game.Modifier = game.Modifier{
	ID:                judgementOfPainID,
	GroupID:           &judgementOfPainID,
	Icon:              "std_chakra_defense",
	Name:              "Judgement of Pain",
	Description:       "Chakra Defense x0.75.",
	ParentDescription: "Other shinobi: Chakra Defense x0.75.",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&judgementOfPainID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.OtherFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stats[game.StatChakraDefense] = game.Round(float64(actor.Stats[game.StatChakraDefense]) * 0.75)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
