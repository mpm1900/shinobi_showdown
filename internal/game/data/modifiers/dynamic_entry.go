package modifiers

import (
	"math/rand/v2"
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var dynamicEntryID = uuid.MustParse("669f5004-666d-4bf0-a071-ac813fdd0be4")
var DynamicEntryTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: dynamicEntryID,
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

			if len(targets) == 0 {
				return transactions
			}

			index := rand.IntN(len(targets))
			target := targets[index]
			mut_ctx := game.MakeContextForActor(target)
			mut_ctx.ModifierID = &dynamicEntryID
			mutation := game.RatioDamage(0.125)
			transaction := game.MakeTransaction(mutation, mut_ctx)
			transactions = append(transactions, transaction)

			return transactions
		},
	},
}

var DynamicEntry game.Modifier = game.Modifier{
	ID:          dynamicEntryID,
	GroupID:     &dynamicEntryID,
	Icon:        "dynamic_entry",
	Name:        "Dynamic Entry",
	Description: "On enter: random enemey loses 1/8th HP.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&dynamicEntryID),
	},
	Triggers: []game.Trigger{
		DynamicEntryTrigger,
	},
}
