package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var entryHazardID = uuid.MustParse("6da2b2f6-f37d-49a5-92d9-49acf561d3d6")
var EntryHazardTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: entryHazardID,
	On:         game.OnActorEnter,
	Check:      game.Match__Player_Player,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
			source, ok := g.GetSource(context)
			if !ok {
				return []game.GameTransaction{}
			}

			mut_ctx := game.MakeContextForActor(source)
			mut := game.RatioDamage(0.0625)
			return []game.GameTransaction{
				game.MakeTransaction(mut, mut_ctx),
			}
		},
	},
}

var EntryHazard = game.Modifier{
	ID:          entryHazardID,
	GroupID:     &entryHazardID,
	Name:        "Entry Hazard",
	Description: "On enter: lose 1/16th HP.",
	Icon:        "entry_hazard",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopPlayer(&entryHazardID),
	},
	Triggers: []game.Trigger{
		EntryHazardTrigger,
	},
}
