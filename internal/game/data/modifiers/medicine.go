package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var medicineID = uuid.MustParse("13778cd8-4ad6-45a1-b769-c605ee0b82af")

var medicineTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: medicineID,
	On:         game.OnStatusAdd,
	Check:      game.Match__TargetActor_SourceActor,
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
			transactions := []game.GameTransaction{}

			if context.ModifierID == nil {
				return transactions
			}

			targets := g.GetTargets(context)
			for _, target := range targets {
				consumeCtx := game.MakeContextForActor(target)
				transactions = append(transactions, game.MakeTransaction(mutations.ConsumeItem, consumeCtx))
				transactions = append(transactions, ClearStatus(g, target)...)
			}

			return transactions
		},
	},
}

var Medicine = game.Modifier{
	ID:          medicineID,
	GroupID:     &medicineID,
	Icon:        "medicine",
	Name:        "Medicine",
	Description: "On status: heal the status, then consume this item.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&medicineID),
	},
	Triggers: []game.Trigger{
		medicineTrigger,
	},
}
