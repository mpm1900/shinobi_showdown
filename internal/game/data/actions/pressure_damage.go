package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var PressureDamage = MakePressureDamage()

func MakePressureDamage() game.Action {
	return makeAttack(AttackConfig{
		ID: uuid.MustParse("22f9cca5-b709-444d-a1d1-72f3707b08cc"),
		Config: makeSpreadAttackConfig(game.ActionConfig{
			Name:        "Pressure Damage",
			Description: "Hits all enemy shinobi. Grants the user Wind nature until end of turn.",
			Nature:      game.Ptr(game.NsWind),
			Accuracy:    game.Ptr(100),
			Power:       game.Ptr(75),
			Stat:        game.Ptr(game.StatChakraAttack),
			Cost:        game.Ptr(30),
			Jutsu:       game.Ninjutsu,
		}),
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_team_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_team_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		BeforeAttack: func(g game.Game, context game.Context, action_config game.ActionConfig) []game.GameTransaction {
			transactions := game.NewTransactionBuilder()

			add_mut := mutations.AddModifiers(false, modifiers.AddNature(game.NsWind, 0))
			add_tx := game.MakeTransaction(add_mut, context)
			transactions.Push(add_tx)

			return transactions.Build()
		},
	})
}
