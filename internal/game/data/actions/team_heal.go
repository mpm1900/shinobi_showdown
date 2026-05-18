package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/modifiers"

	"github.com/google/uuid"
)

var TeamHeal = MakeTeamHeal()

func MakeTeamHeal() game.Action {
	config := makeNoTargetStatusConfig(game.ActionConfig{
		Name:        "Team Heal",
		Description: "All active ally shinobi heal for 1/4th HP.",
		Nature:      game.Ptr(game.NsYang),
		Cost:        game.Ptr(30),
		Jutsu:       game.Senjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("2bb0f69c-fb8a-4390-9041-60444c4a05fc"),
		Config:          config,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(*config.TargetCount),
		Cost:            modifiers.UseStaminaCost(*config.Cost),
		MapContext: func(g game.Game, context game.Context) game.Context {
			actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.TeamFilter))
			for _, t := range actors {
				context.TargetPositionIDs = append(context.TargetPositionIDs, *t.PositionID)
			}
			return context
		},
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				conf, _ := game.GetActiveActionConfig(g, config)
				for _, target := range g.GetTargets(context) {
					heal := game.NewHeal(conf, 0.25)
					transactions = append(
						transactions,
						game.MakeTransaction(heal, game.MakeContextForActor(target)),
					)
				}

				return transactions
			},
		},
	}
}
