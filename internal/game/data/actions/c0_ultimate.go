package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var C0UltimateArt = MakeC0UltimateArt()

func MakeC0UltimateArt() game.Action {
	ID := uuid.MustParse("181d48e6-11d4-45fe-a8a4-09a5fc37c800")

	config := game.ActionConfig{
		Name:        "C0: Ultimate Art",
		Description: "Hits all other active shinobi.",
		Nature:      game.Ptr(game.NsExplosion),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(250),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
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
		OnSuccess: func(g game.Game, _, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			self_dmg := game.RatioDamage(1.0)
			self_dmg_ctx := game.MakeContextForActor(source)
			transactions = append(transactions, game.MakeTransaction(self_dmg, self_dmg_ctx))

			return transactions
		},
	})
}

// proxies
var SelfDestruct = MakeSelfDestruct()

func MakeSelfDestruct() game.Action {
	action := MakeC0UltimateArt()
	action.ID = uuid.MustParse("9cd1049e-d388-47d4-a228-874153cbe5a5")
	action.Config.Name = "Self Destruct"
	action.Config.Stat = game.Ptr(game.StatAttack)
	return action
}
