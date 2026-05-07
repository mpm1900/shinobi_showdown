package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var VacuumBlast = MakeVacuumBlast()

func MakeVacuumBlast() game.Action {
	ID := uuid.MustParse("b5048a55-c3f8-4c80-b70f-447b079ab480")

	config := game.ActionConfig{
		Name:        "Vacuum Blast",
		Description: "Hits all enemy shinobi. Clears weather.",
		Nature:      game.Ptr(game.NsWind),
		Accuracy:    game.Ptr(100),
		Power:       game.Ptr(80),
		Stat:        game.Ptr(game.StatChakraAttack),
		TargetCount: game.Ptr(0),
		Cost:        game.Ptr(30),
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
		OnSuccess: func(g game.Game, context, tcontext game.Context) []game.GameTransaction {
			return modifiers.ClearWeather()
		},
	})
}
