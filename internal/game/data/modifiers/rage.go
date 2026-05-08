package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var rageID = uuid.MustParse("ab1dcd38-82d6-5a88-9123-91caed67da84")
var RageTrigger game.Trigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: rageID,
	On:         game.OnDamageReceive,
	Check:      game.Match__TargetActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			transactions := []game.GameTransaction{}

			targets := g.GetTargets(context)
			for _, target := range targets {
				mut_ctx := game.MakeContextForActor(target)
				mutation := mutations.AddModifiers(false, AttackUpSource)
				transaction := game.MakeTransaction(mutation, mut_ctx)
				transactions = append(transactions, transaction)
			}

			return transactions
		},
	},
}

var Rage game.Modifier = game.Modifier{
	ID:          rageID,
	GroupID:     &rageID,
	Name:        "Rage",
	Icon:        "rage",
	Description: "On damage: gain Attack Up.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&rageID),
	},
	Triggers: []game.Trigger{
		RageTrigger,
	},
}
