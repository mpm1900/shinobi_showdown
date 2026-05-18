package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ShinraTensei = MakeShinraTensei()

func MakeShinraTensei() game.Action {
	config := makeAttackConfig(game.ActionConfig{
		Name:        "Shinra Tensei",
		Description: "Hits all enemy shinobi.",
		Nature:      game.Ptr(game.NsYinYang),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(100),
		Stat:        game.Ptr(game.StatChakraAttack),
		Cost:        game.Ptr(30),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:              uuid.MustParse("fc9fe2ef-32e8-4e20-821e-b98b0a8fb4b7"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		MapContext: func(g game.Game, context game.Context) game.Context {
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))
			for _, t := range other_actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
	})
}
