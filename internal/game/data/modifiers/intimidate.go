package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var intimidateID = uuid.MustParse("0a6d78f6-c09b-5463-ad56-e9d549fc7ca9")
var IntimidateTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: intimidateID,
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
				game.OtherTeamFilter,
			))

			for _, target := range targets {
				mut_ctx := game.MakeContextForActor(target)
				mut_ctx.SourceActorID = context.SourceActorID
				mut_ctx.ModifierID = &intimidateID
				mutation := mutations.AddModifiers(false, AttackDownTarget)
				transaction := game.MakeTransaction(mutation, mut_ctx)
				transactions = append(transactions, transaction)
			}

			return transactions
		},
	},
}

var Intimidate game.Modifier = game.Modifier{
	ID:          intimidateID,
	GroupID:     &intimidateID,
	Icon:        "intimidate",
	Name:        "Intimidate",
	Description: "On enter: all enemies gain Attack Down.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&intimidateID),
	},
	Triggers: []game.Trigger{
		IntimidateTrigger,
	},
}
