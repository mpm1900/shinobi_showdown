package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var shinobiCloakID = uuid.MustParse("1c1ee69e-3758-4ee3-922d-a0f8c1db94ec")

var ShinobiCloak game.Modifier = game.Modifier{
	ID:          shinobiCloakID,
	GroupID:     &shinobiCloakID,
	Icon:        "shinobi_cloak",
	Name:        "Shinobi Cloak",
	Description: "Immune from the secondary effect of actions. (Warded)",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&shinobiCloakID,
			game.MutPriorityPostStagedStats,
			game.ComposeAF(game.SourceFilter, game.ActiveFilter),
			func(g game.Game, actor game.Actor, context game.Context) game.Actor {
				actor.Warded = true
				return actor
			},
		),
	},
	Triggers: []game.Trigger{},
}
