package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var immuneID = uuid.MustParse("21acc3bf-25a1-48b6-8594-5503cc9822b6")
var Immune = game.Modifier{
	ID:          immuneID,
	GroupID:     &immuneID,
	Icon:        "immune",
	Name:        "Immune",
	Description: "Immunity to poison.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		ApplyImmunity(immuneID, Poisoned.ID),
	},
}
