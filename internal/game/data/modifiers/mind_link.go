package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var mindLinkID = uuid.MustParse("9fe5b33a-9d29-4f3d-99ba-ff24eb3f7180")
var MindLink = game.Modifier{
	ID:          mindLinkID,
	Name:        "Mind Link",
	Description: "Protected from damage from allies",
	Icon:        "mind_link",
	Show:        true,
	GroupID:     &mindLinkID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&mindLinkID,
			game.MutPriorityDefault,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter, game.ActiveTargetFilter, game.ActiveTargetTeamSourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.Protected = true
				return a
			},
		),
	},
}
