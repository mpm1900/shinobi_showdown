package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var voiceOfPainID = uuid.MustParse("db8a0316-10be-4f6a-b0c5-a49fb7eb601c")

var VoiceOfPain game.Modifier = game.Modifier{
	ID:                voiceOfPainID,
	GroupID:           &voiceOfPainID,
	Icon:              "std_defense",
	Name:              "Voice of Pain",
	Description:       "Defense x0.75.",
	ParentDescription: "Other shinobi: Defense x0.75.",
	Show:              true,
	Duration:          game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&voiceOfPainID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.OtherFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Stats[game.StatDefense] = game.Round(float64(actor.Stats[game.StatDefense]) * 0.75)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
