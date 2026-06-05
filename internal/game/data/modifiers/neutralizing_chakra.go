package modifiers

import (
	"shinobi_showdown/internal/game"
	"shinobi_showdown/internal/game/data/mutations"
	"slices"

	"github.com/google/uuid"
)

/**
 * This modifier is a bit complex so here's a commet to explain things.
 * This modifier has to use triggers to assign an empty "AuxAbility" so that
 * we can actually modifier what the modifier sources are.
 *
 * Then, another limitation of the current system is an actor's modifiers are not active
 * at the time of OnActorLeave trigger resolution. This means that all other actors must be
 * given their own OnActorLeave triggers that all listen for the source actor's trigger.
 */

var NoopAbility game.Modifier = game.Modifier{
	ID:             uuid.New(),
	Name:           "",
	Show:           false,
	Duration:       game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{},
	Triggers:       []game.Trigger{},
}

var neutralizingChakraID = uuid.MustParse("ca02f87f-a32c-4977-bc75-5720fffc0475")

var NeutralizingChakra game.Modifier = game.Modifier{
	ID:             neutralizingChakraID,
	GroupID:        &neutralizingChakraID,
	Name:           "Neutralizing Chakra",
	Icon:           "neutralizing_chakra",
	Description:    "On enter: Disable all active abilities",
	Show:           true,
	Duration:       game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{},
	Triggers: []game.Trigger{
		// source entered, wipe everyones ability.
		{
			ID:         uuid.New(),
			ModifierID: neutralizingChakraID,
			On:         game.OnActorEnter,
			Check:      game.Match__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP1,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					targets := g.GetActorsFilters(context, game.ComposeAF(
						game.ActiveFilter,
						game.AliveFilter,
						game.OtherFilter,
					))
					for _, target := range targets {
						mut_ctx := context
						mut_ctx.ModifierID = &neutralizingChakraID
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

						mod_ctx := game.MakeContextForActor(target)
						mod_ctx.TargetActorIDs = []uuid.UUID{source.ID}
						mod := mutations.AddModifiers(false, returnAbility)
						mod_tx := game.MakeTransaction(mod, mod_ctx)
						transactions = append(transactions, mod_tx)
					}

					return transactions
				},
			},
		},
		// new actor just entered, remove it's ability trigger
		{
			ID:         uuid.New(),
			ModifierID: neutralizingChakraID,
			On:         game.OnActorEnter,
			Check:      game.NotMatch__SourceActor_SourceActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP1,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					mut_ctx := context
					mut_ctx.ModifierID = &neutralizingChakraID
					mut_ctx.TargetActorIDs = []uuid.UUID{source.ID}
					mutation := game.GameMutation{
						Delta: func(p, g game.Game, context game.Context) game.Game {
							g.UpdateActor(source.ID, func(a game.Actor) game.Actor {
								a.AuxAbility = &NoopAbility
								return a
							})
							return g
						},
					}
					transaction := game.MakeTransaction(mutation, mut_ctx)
					transactions = append(transactions, transaction)

					neutralizers := g.GetActorsFilters(context, game.ComposeAF(
						game.ActiveFilter,
						game.AliveFilter,
						game.OtherFilter,
						game.HasAppliedModifier(neutralizingChakraID),
					))
					for _, neutralizer := range neutralizers {
						mod_ctx := game.MakeContextForActor(source)
						mod_ctx.TargetActorIDs = []uuid.UUID{neutralizer.ID}
						mod_tx := game.MakeTransaction(mutations.AddModifiers(false, returnAbility), mod_ctx)
						transactions = append(transactions, mod_tx)
					}

					return transactions
				},
			},
		},
	},
}

var returnAbilityID = uuid.MustParse("e10a810c-60b9-4d25-afdf-22404f946848")
var returnAbility = game.Modifier{
	ID:             returnAbilityID,
	GroupID:        &returnAbilityID,
	Name:           "Return Ability",
	Show:           false,
	Duration:       game.ModifierDurationInf,
	ActorMutations: []game.ActorMutation{},
	Triggers: []game.Trigger{
		{
			ID:         uuid.New(),
			ModifierID: returnAbilityID,
			On:         game.OnActorLeave,
			Check:      game.Match__SourceActor_TargetActor,
			ActionMutation: game.ActionMutation{
				Priority: game.ActionPriorityP1,
				Filter:   game.TrueGameFilter,
				Delta: func(p, g game.Game, context game.Context) []game.GameTransaction {
					transactions := []game.GameTransaction{}
					source, ok := g.GetSource(context)
					if !ok {
						return transactions
					}

					remove_mut := mutations.RemoveModifierWhere(func(tx game.Transaction[game.Modifier]) bool {
						return tx.Mutation.ID == returnAbilityID && slices.Contains(tx.Context.TargetActorIDs, source.ID)
					})
					transactions = append(transactions, game.MakeTransaction(remove_mut, context))

					neutralizers := g.GetActorsFilters(context, game.ComposeAF(
						game.ActiveFilter,
						game.AliveFilter,
						game.OtherFilter,
						game.HasAppliedModifier(neutralizingChakraID),
					))
					if len(neutralizers) > 0 {
						return transactions
					}

					targets := g.GetActorsFilters(context, game.ComposeAF(
						game.ActiveFilter,
						game.AliveFilter,
						game.OtherFilter,
					))

					for _, target := range targets {
						mut_ctx := context
						mut_ctx.ModifierID = &neutralizingChakraID
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
