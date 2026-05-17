package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var shinobiVestID = uuid.MustParse("a05df03d-bfc7-5a0a-a01b-831fa1deff3b")

var ShinobiVest game.Modifier = game.Modifier{
	ID:          shinobiVestID,
	GroupID:     &shinobiVestID,
	Icon:        "shinobi_vest",
	Name:        "Shinobi Vest",
	Description: "Chakra Defense x1.5; attacking actions only.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&shinobiVestID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				for i := range actor.Actions {
					actor.Actions[i].State.Disabled = !actor.Actions[i].Meta.Switch && actor.Actions[i].Config.Power == nil
				}

				actor.Stats[game.StatChakraDefense] = game.Round(float64(actor.Stats[game.StatChakraDefense]) * 1.5)
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
