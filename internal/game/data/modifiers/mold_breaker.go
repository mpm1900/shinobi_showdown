package modifiers

import (
	"shinobi_showdown/internal/game"

	"github.com/google/uuid"
)

var moldBreakerID = uuid.MustParse("334a4765-5b63-4164-8c1a-2aa070de7754")

var MoldBreaker game.Modifier = game.Modifier{
	ID:             moldBreakerID,
	GroupID:        &moldBreakerID,
	Name:           "Mold Breaker",
	Icon:           "mold_breaker",
	Description:    "Disabled abilities while targeted.",
	Show:           true,
	Duration:       game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: moldBreakerID,
			On:         game.OnActionStart,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP1,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					targets := g.GetActorsFilters(context, game.ComposeAF(
						game.TargetFilter,
					))
					for _, target := range targets {
						mut_ctx := context
						mut_ctx.ModifierID = &moldBreakerID
						mut_ctx.TargetActorIDs = []uuid.UUID{target.ID}
						mutation := game.GameMutation{
							Delta: func(p, g game.Game, context game.Context) game.Game {
								g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
									a.AuxAbility = &NoopAbility
									return a
								})
								return g
							},
						}
						transaction := game.MakeTransaction(mutation, mut_ctx)
						transactions = append(transactions, transaction)
					}

					return transactions
				},
			},
		},
		{
			ID:         uuid.New(),
			ModifierID: moldBreakerID,
			On:         game.OnActionEnd,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP1,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					targets := g.GetActorsFilters(context, game.ComposeAF(
						game.TargetFilter,
					))
					for _, target := range targets {
						mut_ctx := context
						mut_ctx.ModifierID = &moldBreakerID
						mut_ctx.TargetActorIDs = []uuid.UUID{target.ID}
						mutation := game.GameMutation{
							Delta: func(p, g game.Game, context game.Context) game.Game {
								g.UpdateActor(target.ID, func(a game.Actor) game.Actor {
									a.AuxAbility = nil
									return a
								})
								return g
							},
						}
						transaction := game.MakeTransaction(mutation, mut_ctx)
						transactions = append(transactions, transaction)
					}

					return transactions
				},
			},
		},
	},
}
