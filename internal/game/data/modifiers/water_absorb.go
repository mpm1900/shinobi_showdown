package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var waterAbsorbID = uuid.MustParse("ebc37353-4eff-5b89-b7f5-361849f2028b")

var WaterAbsorb game.Modifier = game.Modifier{
	ID:          waterAbsorbID,
	Name:        "Water Absorb",
	Icon:        "water_absorb",
	Description: "Water damage heals this shinobi.",
	Show:        true,
	GroupID:     &waterAbsorbID,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		{
			ModifierGroupID: &waterAbsorbID,
			Mutation: game.Mutation[game.Game, game.Actor, game.Actor]{
				Priority: game.MutPrioritySet,
				Filter:   game.SourceFilter,
				Delta: func(p game.Game, a game.Actor, c game.Context) game.Actor {
					action, ok := p.GetActiveAction()
					if !ok || action.Config.Nature == nil {
						return a
					}

					if *action.Config.Nature != game.NsWater {
						return a
					}

					a.NatureResistance[game.NatureWater] *= -1.0
					return a
				},
			},
		},
	},
}
