package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var Earthquake = MakeEarthquake()

func MakeEarthquake() game.Action {
	config := makeSpreadAttackConfig(game.ActionConfig{
		Name:        "Earthquake",
		Description: "Hits all other grounded shinobi. Applies Rocky terrain.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	})

	done := false // this closure is really stupid but works to define a "once per attack on success"
	action := makeAttack(AttackConfig{
		ID:              uuid.MustParse("80e197af-299c-448d-a431-12473fa13866"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherFilter, game.RHasState(game.ActorStateGrounded)))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		BeforeAttack: func(g game.Game, ctx game.Context, action_config game.ActionConfig) []game.GameTransaction {
			done = false
			return []game.GameTransaction{}
		},
		OnSuccess: func(g game.Game, context, tctx game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			if done {
				return transactions.Build()
			}

			transactions.Push(modifiers.ApplyTerrain(g, context, game.GameTerrainRocky, modifiers.RockyTerrain()))

			done = true
			return transactions.Build()
		},
	})

	return action
}
