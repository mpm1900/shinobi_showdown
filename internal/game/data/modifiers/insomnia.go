package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var insomniaID = uuid.MustParse("21acc3bf-25a1-48b6-8594-5503cc9822b6")
var Insomnia = game.Modifier{
	ID:          insomniaID,
	GroupID:     &insomniaID,
	Icon:        "insomnia",
	Name:        "Insomnia",
	Description: "Immunity to sleep.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		ApplyImmunity(insomniaID, Sleeping.ID),
	},
}
