package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var GreatFireAnnihilation = MakeGreatFireAnnihilation()

func MakeGreatFireAnnihilation() game.Action {
	ID := uuid.MustParse("d97ee3bb-7afa-47de-9f8d-2ee77ba6dfe6")

	config := game.ActionConfig{
		Name:        "Great Fire Annihilation",
		Description: "Hits all enemy shinobi. 20% chance to burn targets.",
		Nature:      game.Ptr(game.NsFire),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(90),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		Cost:        game.Ptr(30),
		Cooldown:    game.Ptr(2),
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
				transactions = append(transactions, modifiers.ChanceBurn(config, g, context, target, 20)...)
			}

			return transactions
		},
	})
}
