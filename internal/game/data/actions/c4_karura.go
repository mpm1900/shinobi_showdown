package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var c4KaruraModifierID = uuid.MustParse("63ecfb63-fcc9-4240-9f03-525fdc04c6d0")
var c4KaruraID = uuid.MustParse("4d1f216c-83b9-468d-86e5-595dcc8d32c7")

var c4KaruraTrigger = game.Trigger{
	ID:         uuid.New(),
	ModifierID: c4KaruraModifierID,
	On:         game.OnTurnEnd,
	Check: func(p, g game.Game, context game.Context, tx game.Transaction[game.Modifier]) bool {
		return true
	},
	ActionMutation: game.ActionMutation{
		Priority: game.ActionPriorityDefault,
		Filter:   game.TrueGameFilter,
		Delta: func(p game.Game, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
			transactions := []game.GameTransaction{}

			if context.SourceActorID == nil {
				return transactions
			}

			mut := game.RatioDamage(0.5)
			other_actors := g.GetActorsFilters(context, game.ComposeAF(game.ActiveFilter, game.OtherTeamFilter))

			for _, target := range other_actors {
				mut_ctx := game.MakeContextForActor(target).WithSource(*context.SourceActorID)
				transactions = append(transactions, game.MakeTransaction(mut, mut_ctx))
			}

			return transactions
		},
	},
}

var c4KaruraModifier = game.Modifier{
	ID:          c4KaruraModifierID,
	GroupID:     &c4KaruraModifierID,
	Name:        "C4: Karura",
	Description: "On turn end: each enemy shinobi loses half of their HP.",
	Icon:        "c4_karura",
	Show:        false,
	Duration:    3,
	Delay:       2,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&c4KaruraModifierID),
	},
	Triggers: []game.Trigger{
		c4KaruraTrigger,
	},
}

var C4Karura = MakeC4Karura()

func MakeC4Karura() game.Action {
	config := game.ActionConfig{
		Name:        "C4: Karura",
		Nature:      game.Ptr(game.NsExplosion),
		Jutsu:       game.Ninjutsu,
		Cost:        game.Ptr(80),
		Description: "In 2 turns, each enemy shinobi loses half of their HP. This effect is nullified if the user dies or switches out.",
	}
	return game.Action{
		ID:              c4KaruraID,
		Config:          config,
		TargetType:      game.TargetActorID,
		TargetPredicate: game.NoneFilter,
		ContextValidate: game.TargetLengthFilter(0),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter:   game.SourceIsAlive,
			Delta: func(p game.Game, g game.Game, context game.Context) []game.GameTransaction {
				transactions := []game.GameTransaction{}

				mutation := mutations.AddModifiers(false, c4KaruraModifier)
				transaction := game.MakeTransaction(mutation, context)
				transactions = append(transactions, transaction)

				return transactions
			},
		},
	}
}
