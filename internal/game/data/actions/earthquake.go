package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Earthquake = MakeEarthquake()

func MakeEarthquake() game.Action {
	ID := uuid.MustParse("80e197af-299c-448d-a431-12473fa13866")

	config := game.ActionConfig{
		Name:        "Earthquake",
		Description: "Hits all other active shinobi. Applies Rocky terrain.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		TargetCount: game.Ptr(0),
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
				modifiers.ApplyTerrain(g, context, game.GameTerrainRocky, modifiers.RockyTerrain())...,
			)

			done = true
			return transactions
		},
	})

	return action
}
