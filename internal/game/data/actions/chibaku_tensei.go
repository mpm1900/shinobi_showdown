package actions

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var ChibakuTensei = MakeChibakuTensei()

func MakeChibakuTensei() game.Action {
	ID := uuid.MustParse("a2f1123f-4fef-4f5e-bd8c-f17d7bb3b43a")

	config := game.ActionConfig{
		Name:        "Chibaku Tensei",
		Description: "Hits all enemy shinobi.",
		Nature:      game.Ptr(game.NsEarth),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(110),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		TargetType:  game.TargetPositionID,
		Cost:        game.Ptr(100),
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
	})
}
