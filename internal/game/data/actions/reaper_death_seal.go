package actions

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"

	"github.com/google/uuid"
)

var ReaperDeathSeal = MakeReaperDeathSeal()

func MakeReaperDeathSeal() game.Action {
	config := makeStatusConfig(game.ActionConfig{
		Name:        "Reaper Death Seal",
		Description: "Bonds user and target. When either bonded shinobi dies, both die.",
		Nature:      game.Ptr(game.NsYin),
		Jutsu:       game.Fuinjutsu,
	})

	return game.Action{
		ID:              uuid.MustParse("0cf47657-82bf-4a82-b933-cd7a762e0327"),
		Config:          config,
		TargetPredicate: game.ComposeAF(game.OtherFilter, game.TargetableFilter),
		ContextValidate: game.PositionsLengthFilter(*config.TargetCount),
		ActionMutation: game.ActionMutation{
			Priority: game.ActionPriorityDefault,
			Filter: game.ComposeGF(
				game.SourceIsAlive,
			),
			Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
				transactions := []game.GameTransaction{}

				source, ok := g.GetSource(context)
				if !ok {
					return transactions
				}

				for _, target := range g.GetTargets(context) {
					source_ctx := game.MakeContextForActor(source)
					source_ctx.TargetActorIDs = []uuid.UUID{target.ID}
					source_ctx.ModifierID = &SourceBonded.ID
					source_mut := mutations.AddModifiers(false, SourceBonded)
					source_tx := game.MakeTransaction(source_mut, source_ctx)
					transactions = append(transactions, source_tx)

					target_ctx := game.MakeContextForActor(target)
					target_ctx.TargetActorIDs = []uuid.UUID{source.ID}
					target_ctx.ModifierID = &targetBondedID
					target_mut := mutations.AddModifiers(true, TargetBonded)
					target_tx := game.MakeTransaction(target_mut, target_ctx)
					transactions = append(transactions, target_tx)
				}

				return transactions
			},
		},
	}
}

var sourceBondedID = uuid.New()
var SourceBonded game.Modifier = game.Modifier{
	ID:          sourceBondedID,
	GroupID:     &sourceBondedID,
	Name:        "Bonded",
	Description: "When a bonded shinobi dies, so does the other.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&sourceBondedID),
	},
	Triggers: []game.Trigger{
		// 1. Death trigger
		{
			ID:         uuid.New(),
			ModifierID: sourceBondedID,
			On:         game.OnDeath,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					transactions := []game.GameTransaction{}
					if context.SourceActorID == nil {
						return transactions
					}

					var bonded_partner_ID *uuid.UUID = nil
					for _, tx := range g.Modifiers {
						if tx.Mutation.ID == sourceBondedID && tx.Context.SourceActorID != nil && *tx.Context.SourceActorID == *context.SourceActorID {
							if len(tx.Context.TargetActorIDs) > 0 {
								bonded_partner_ID = &tx.Context.TargetActorIDs[0]
							}
							break
						}
					}

					if bonded_partner_ID != nil {
						partner_ctx := game.NewContext().WithSource(*context.SourceActorID).WithTargetIDs([]uuid.UUID{*bonded_partner_ID})
						mut := game.RatioDamage(1.0)
						tx := game.MakeTransaction(mut, partner_ctx)
						transactions = append(transactions, tx)
					}

					return transactions
				},
			},
		},
		// 2. Break trigger
		{
			ID:         uuid.New(),
			ModifierID: sourceBondedID,
			On:         game.OnActorLeave,
			Check:      game.Match__SourceActor_TargetActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					transactions := []game.GameTransaction{}
					if context.SourceActorID == nil {
						return transactions
					}

					mut := mutations.RemoveModifierWhere(func(tx game.Transaction[game.Modifier]) bool {
						for _, tid := range tx.Context.TargetActorIDs {
							if tx.Mutation.ID == sourceBondedID && tid == *context.SourceActorID {
								return true
							}
						}

						return false
					})
					tx := game.MakeTransaction(mut, game.NewContext())
					transactions = append(transactions, tx)

					return transactions
				},
			},
		},
	},
}

var targetBondedID = uuid.New()
var TargetBonded game.Modifier = game.Modifier{
	ID:          targetBondedID,
	GroupID:     &targetBondedID,
	Name:        "Bonded",
	Description: "When a bonded shinobi dies, so does the other.",
	Show:        true,
	Duration:    game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{
		game.NewNoopSource(&targetBondedID),
	},
	Triggers: []game.Trigger{
		// 1. Death trigger
		{
			ID:         uuid.New(),
			ModifierID: targetBondedID,
			On:         game.OnDeath,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					transactions := []game.GameTransaction{}
					if context.SourceActorID == nil {
						return transactions
					}

					var bonded_partner_ID *uuid.UUID = nil
					for _, tx := range g.Modifiers {
						if tx.Mutation.ID == targetBondedID && tx.Context.SourceActorID != nil && *tx.Context.SourceActorID == *context.SourceActorID {
							if len(tx.Context.TargetActorIDs) > 0 {
								bonded_partner_ID = &tx.Context.TargetActorIDs[0]
							}
							break
						}
					}

					if bonded_partner_ID != nil {
						partner_ctx := game.NewContext().WithSource(*context.SourceActorID).WithTargetIDs([]uuid.UUID{*bonded_partner_ID})
						mut := game.RatioDamage(1.0)
						tx := game.MakeTransaction(mut, partner_ctx)
						transactions = append(transactions, tx)
					}

					return transactions
				},
			},
		},
		// 2. Break trigger
		{
			ID:         uuid.New(),
			ModifierID: targetBondedID,
			On:         game.OnActorLeave,
			Check:      game.Match__SourceActor_TargetActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityDefault,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.Transaction[game.GameMutation] {
					transactions := []game.Transaction[game.GameMutation]{}
					if context.SourceActorID == nil {
						return transactions
					}

					mut := mutations.RemoveModifierWhere(func(tx game.Transaction[game.Modifier]) bool {
						for _, tid := range tx.Context.TargetActorIDs {
							if tx.Mutation.ID == targetBondedID && tid == *context.SourceActorID {
								return true
							}
						}

						return false
					})
					tx := game.MakeTransaction(mut, game.NewContext())
					transactions = append(transactions, tx)

					return transactions
				},
			},
		},
	},
}
