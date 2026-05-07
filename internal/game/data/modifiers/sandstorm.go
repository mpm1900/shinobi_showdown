package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var sandAuraID = uuid.MustParse("b16bfb0c-8131-4522-8f97-7c5775e7df05")

var SandAuraTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: sandAuraID,
	On:         game.OnActorEnter,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			return ApplySandstorm(g, context)
		},
	},
}

var SandAura = game.Modifier{
	ID:          sandAuraID,
	GroupID:     &sandAuraID,
	Icon:        "sand_aura",
	Name:        "Sand Aura",
	Description: "On enter: start sandstorm.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sandAuraID),
	},
	Triggers: []game.Trigger{
		SandAuraTrigger,
	},
}
