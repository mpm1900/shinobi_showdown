package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var natureSpecialistID = uuid.MustParse("1b7fa30d-e6c8-43b4-bd4a-10f2f754b9cd")
var NatureSpecialist = game.Modifier{
	ID:          natureSpecialistID,
	Name:        "Nature Specialist",
	Description: "User's STAB multiplier is 2x instead of 1.5x.",
	Icon:        "nature_specialist",
	Show:        true,
	GroupID:     &natureSpecialistID,
	Duration:    0,
	ActorMutations: []game.ActorMutation{
		game.MakeActorMutation(
			&natureSpecialistID,
			game.MutPrioritySet,
			game.ComposeAF(game.ActiveFilter, game.SourceFilter),
			func(g game.Game, a game.Actor, c game.Context) game.Actor {
				a.StabMultiplier = 2
				return a
			},
		),
	},
}
