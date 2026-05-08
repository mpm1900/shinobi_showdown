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
		Cost:        game.Ptr(10),
		Power:       game.Ptr(20),
		Accuracy:    game.Ptr(100),
		TargetCount: game.Ptr(0),
		Stat:        game.Ptr(game.StatChakraAttack),
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	done := false // this closure is really stupid but works to define a "once per attack on success"

	return makeAttack(AttackConfig{
		ID:              burningAshID,
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		BeforeAttack: func(g game.Game, ctx game.Context) []game.GameTransaction {
			done = false
			return []game.GameTransaction{}
		},
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 20)...)
			}

			if done {
				return transactions
			}

			state, _ := g.GetState(context)
			if state.Terrain == game.GameTerrainFlamable {
				return transactions
			}

			filter := modifiers.FilterTerrain()
			transactions = append(transactions, filter)

			mod := modifiers.FlamableTerrain()
			mod.Duration = 4
			mut := mutations.AddModifiers(false, mod)
			terrain_tx := game.MakeTransaction(mut, game.NewContext().WithSource(*context.SourceActorID))
			transactions = append(transactions, terrain_tx)

			done = true
			return transactions
		},
	})
}
