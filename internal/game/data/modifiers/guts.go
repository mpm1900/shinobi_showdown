package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var gutsID = uuid.MustParse("efda299f-bedd-4c3a-8e67-146a994a231d")

func makeGutsMutation() game.ActorMutation {
	mut := game.NewNoopSource(&gutsID)
	mut.Priority = game.MutPriorityImmunity
	return mut
}

var Guts game.Modifier = game.Modifier{
	ID:          gutsID,
	GroupID:     &gutsID,
	Icon:        "guts",
	Name:        "Guts",
	Description: "Status conditions increase Attack by 50%.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		makeGutsMutation(),
	},
}
