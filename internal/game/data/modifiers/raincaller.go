package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var raincallerID = uuid.MustParse("912e5e72-263e-5f6f-8f9f-32d1a746cc49")

var RaincallerTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: raincallerID,
	On:         game.OnActorEnter,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			return ApplyRain(g, context)
		},
	},
}

var Raincaller game.Modifier = game.Modifier{
	ID:          raincallerID,
	GroupID:     &raincallerID,
	Icon:        "raincaller",
	Name:        "Raincaller",
	Description: "On enter: start rain.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&raincallerID),
	},
	Triggers: []game.Trigger{
		RaincallerTrigger,
	},
}
