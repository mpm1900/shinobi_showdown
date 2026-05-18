package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var DarkSwamp = MakeDarkSwamp()

func MakeDarkSwamp() game.Action {
	config := makeSpreadAttackConfig(game.ActionConfig{
		Name:        "Dark Swamp",
		Description: "Hits all enemy shinobi. 30% chance to stun targets.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(75),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(100),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:              uuid.MustParse("82947ffd-47a1-4bca-be69-7e545f12e0ab"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
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
				transactions = append(transactions, modifiers.ChanceModifier(config, g, context, target, modifiers.Stunned, 30)...)
			}

			return transactions
		},
	})
}
