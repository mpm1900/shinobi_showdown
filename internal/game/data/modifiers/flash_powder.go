package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var flashPowderID = uuid.MustParse("389b4222-7dd2-4b5f-9876-1b154505a798")

var FlashPowder = game.Modifier{
	ID:          flashPowderID,
	Name:        "Flash Powder",
	Icon:        "flash_powder",
	Description: "Shinobi targeting the holder have 10% reduced accuracy",
	Show:        true,
	GroupID:     &flashPowderID,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&flashPowderID),
		{
			ModifierGroupID: &waterAbsorbID,
			Mutation: game.Mutation[game.Game, game.Actor, game.Actor]{
				Priority: game.MutPrioritySet,
				Filter:   game.ComposeAF(game.OtherFilter, game.ActiveTransactionFilter),
				Delta: func(g game.Game, a game.Actor, context game.Context) game.Actor {
					tx := g.ActiveTransaction
					if tx == nil {
						return a
					}

					if context.SourceActorID == nil {
						return a
					}

					found := false
					targets := g.GetTargets(tx.Context)
					for _, target := range targets {
						if target.ID == *context.SourceActorID {
							found = true
							break
						}
					}

					if found {
						a.ActionAccuracyOffset -= 10
					}

					return a
				},
			},
		},
	},
}
