package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var SearingMigraine = MakeSearingMigraine()

func MakeSearingMigraine() game.Action {
	config := makeSpreadAttackConfig(game.ActionConfig{
		Name:        "Searing Migraine",
		Description: "Hits all enemy shinobi. Grants the user Fire nature until end of turn.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(75),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:              uuid.MustParse("dc6edab6-535f-508f-b791-e197283eae86"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_team_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_team_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		BeforeAttack: func(g game.Game, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			var transactions = []game.GameTransaction{}
			add_mut := mutations.AddModifiers(false, modifiers.AddNature(game.NsFire, 0))
			add_tx := game.MakeTransaction(add_mut, context)
			transactions = append(transactions, add_tx)
			return transactions
		},
	})
}
