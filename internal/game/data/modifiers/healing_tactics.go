package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var healingTacticsID = uuid.MustParse("0b932a5b-0c00-4fda-9629-063c936c2f91")
var HealingTacticsTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: healingTacticsID,
	On:         game.OnActorEnter,
	Check:      game.Match__SourceActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}
			targets := g.GetActorsFilters(context, game.ComposeAF(
				game.ActiveFilter,
				game.AliveFilter,
				game.OtherFilter,
				game.TeamFilter,
			))

			for _, target := range targets {
				mut_ctx := game.MakeContextForActor(target)
				mut_ctx.ModifierID = &healingTacticsID
				mutation := game.RatioHeal(0.25)
				transaction := game.MakeTransaction(mutation, mut_ctx)
				transactions = append(transactions, transaction)
			}

			return transactions
		},
	},
}

var HealingTactics game.Modifier = game.Modifier{
	ID:          healingTacticsID,
	GroupID:     &healingTacticsID,
	Icon:        "healing_tactics",
	Name:        "Healing Tactics",
	Description: "On enter: allies heal 1/4th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&healingTacticsID),
	},
	Triggers: []game.Trigger{
		HealingTacticsTrigger,
	},
}
