package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var burningAshID = uuid.MustParse("8182c2a7-48ea-4d71-ae05-641988c70993")
var BurningAsh = MakeBurningAsh()

func MakeBurningAsh() game.Action {
	config := game.ActionConfig{
		Name:        "Burning Ash",
		Nature:      game.Ptr(game.NsFire),
		Jutsu:       game.Ninjutsu,
		Description: "Sets Flamable Terrain.",
	}

	return game.Action{
		ID:              burningAshID,
		Config:          config,
		TargetType:      game.TargetPositionID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				state, _ := g.GetState(context)
				if state.Terrain == game.GameTerrainFlamable {
					return transactions
				}

				filter := modifiers.FilterTerrain()
				transactions = append(transactions, filter)

				mod := modifiers.FloodedTerrain()
				mod.Duration = 4
				mut := mutations.AddModifiers(false, mod)
				terrain_tx := game.MakeTransaction(mut, game.NewContext().WithSource(*context.SourceActorID))
				transactions = append(transactions, terrain_tx)

				return transactions
			},
		},
	}
}
