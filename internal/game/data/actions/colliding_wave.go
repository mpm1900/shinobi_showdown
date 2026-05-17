package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var CollidingWave = MakeCollidingWave()

func MakeCollidingWave() game.Action {
	ID := uuid.MustParse("74d5a7d7-cb62-58b4-9ace-e80bf7f0fd40")

	config := game.ActionConfig{
		Name:        "Colliding Wave",
		Description: "Hits all other active shinobi. Sets flooded terrain.",
		Nature:      game.Ptr(game.NsWater),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	done := false // this closure is really stupid but works to define a "once per attack on success"
	action := makeAttack(AttackConfig{
		ID:              ID,
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		BeforeAttack: func(g game.Game, ctx game.Context) []game.GameTransaction {
			done = false
			return []game.GameTransaction{}
		},
		OnSuccess: func(g game.Game, context, tctx game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			if done {
				return transactions
			}

			transactions = append(
				transactions,
				modifiers.ApplyTerrain(g, context, game.GameTerrainFlooded, modifiers.FloodedTerrain())...,
			)

			done = true
			return transactions
		},
	})

	return action
}
