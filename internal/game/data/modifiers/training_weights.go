package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var trainingWeightsID = uuid.MustParse("c3dc0104-af0f-46c2-828d-8272f4d0db8b")
var TrainingWeights game.Modifier = game.Modifier{
	ID:          trainingWeightsID,
	GroupID:     &trainingWeightsID,
	Icon:        "training_weights",
	Name:        "Training Weights",
	Description: "Holder moves last in speed priority bracket.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&trainingWeightsID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				for i, _ := range actor.Actions {
					actor.Actions[i].Config.SubPriority = -1
				}
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
