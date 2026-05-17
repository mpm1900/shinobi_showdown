package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var DarkSwamp = MakeDarkSwamp()

func MakeDarkSwamp() game.Action {
	ID := uuid.MustParse("82947ffd-47a1-4bca-be69-7e545f12e0ab")

	config := game.ActionConfig{
		Name:        "Dark Swamp",
		Description: "Hits all enemy shinobi. 30% chance to stun targets.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(90),
		Power:       game.Ptr(75),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(100),
		Cooldown:    game.Ptr(0),
		Jutsu:       game.Ninjutsu,
		CritChance:  game.Ptr(getCriticalStage(0)),
		CritMod:     1.5,
	}

	return makeAttack(AttackConfig{
		ID:              ID,
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
