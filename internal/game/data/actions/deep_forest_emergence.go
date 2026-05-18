package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var DeepForestEmergence = MakeDeepForestEmergence()

func MakeDeepForestEmergence() game.Action {
	config := makeSpreadAttackConfig(game.ActionConfig{
		Name:        "Deep Forest Emergence",
		Description: "Hits all enemy shinobi.",
		Nature:      game.Ptr(game.NsWood),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(110),
		Stat:        game.Ptr(game.StatAttack),
		Cost:        game.Ptr(100),
		Cooldown:    game.Ptr(2),
		Jutsu:       game.Ninjutsu,
	})

	return makeAttack(AttackConfig{
		ID:              uuid.MustParse("93c64751-5625-4887-8d06-90777403e0a1"),
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
