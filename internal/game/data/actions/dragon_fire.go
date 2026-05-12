package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var DragonFire = MakeDragonFire()

func MakeDragonFire() game.Action {
	ID := uuid.MustParse("dca159df-75fb-4cf0-85c5-db69d987a029")

	config := game.ActionConfig{
		Name:        "Dragon Fire",
		Description: "25% chance to burn target. In Flammable Terrain, hits all enemies.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(1),
		Cost:        game.Ptr(0),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:     ID,
		Config: config,
		MapContext: func(g game.Game, context game.Context) game.Context {
			state, _ := g.GetState(context)
			if state.Terrain != game.GameTerrainFlamable {
				return context
			}

			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 25)...)
			}

			return transactions
		},
	})
}
