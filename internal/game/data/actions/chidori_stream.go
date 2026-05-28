package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var ChidoriStream = MakeChidoriStream()

func MakeChidoriStream() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("8973e59d-dc43-4eac-99e2-68e6ca114aa9"),
		Config: makeSpreadAttackConfig(game.ActionConfig{
			Name:        "Chidori Stream",
			Description: "Hits all enemy shinobi. 10% chance to paralyze targets.",
			Nature:      game.Ptr(game.NsLightning),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(80),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(30),
			Jutsu:       game.Ninjutsu,
		}),
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		OnSuccess: func(g game.Game, _, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			targets := g.GetTargets(context)
			for _, target := range targets {
				transactions.Push(modifiers.ChanceParalysis(action_config, g, context, target, 10))
			}

			return transactions.Build()
		},
	})
}
