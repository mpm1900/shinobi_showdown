package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var C0UltimateArt = MakeC0UltimateArt()

func MakeC0UltimateArt() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("181d48e6-11d4-45fe-a8a4-09a5fc37c800"),
		Config: makeAttackConfig(game.ActionConfig{
			Name:        "C0: Ultimate Art",
			Description: "Hits all other active shinobi. User dies.",
			Nature:      game.Ptr(game.NsExplosion),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(250),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(80),
			Jutsu:       game.Ninjutsu,
		}),
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()
			source, ok := g.GetSource(context)
			if !ok {
				return transactions.Build()
			}

			self_dmg := mutations.KillSource()
			self_dmg_ctx := game.MakeContextForActor(source)
			transactions.Push(game.MakeTransaction(self_dmg, self_dmg_ctx))

			return transactions.Build()
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
