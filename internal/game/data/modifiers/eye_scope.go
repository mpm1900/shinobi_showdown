package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var eyeScopeID = uuid.MustParse("33b874ad-0b1e-4862-ab94-5d7d3b975863")
var EyeScope game.Modifier = game.Modifier{
	ID:          eyeScopeID,
	GroupID:     &eyeScopeID,
	Icon:        "eye_scope",
	Name:        "Eye Scope",
	Description: "+15% critical chance.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&eyeScopeID,
			game.MutPriorityDefault,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				for i, action := range actor.Actions {
					if action.Config.CritChance != nil {
						actor.Actions[i].Config.CritChance = game.Ptr(*action.Config.CritChance + 15)
					}
				}
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
