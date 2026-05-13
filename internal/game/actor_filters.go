package game

import (
	"slices"

	"github.com/google/uuid"
)

/**
 * Actor Filters
 */

type ActorFilter func(Game, Actor, Context) bool

func ComposeAF(filters ...ActorFilter) ActorFilter {
	return func(game Game, actor Actor, context Context) bool {
		for _, filter := range filters {
			if !filter(game, actor, context) {
				return false
			}
		}

		return true
	}
}
func AllFilter(game Game, actor Actor, context Context) bool {
	return true
}
func NoneFilter(game Game, actor Actor, context Context) bool {
	return false
}
func OtherFilter(game Game, actor Actor, context Context) bool {
	if context.SourceActorID == nil {
		return false
	}
	return actor.ID != *context.SourceActorID
}

/**
 * This filter doesn't need to be a resolved filter
 * actor.Alive is a special case were modifiers don't modify it, we just mutate it
 * so this is a safe check to make without resoloving the Actor to a ResolvedActor
 */
func AliveFilter(game Game, actor Actor, context Context) bool {
	return actor.Alive
}
func NotAliveFilter(game Game, actor Actor, context Context) bool {
	return !actor.Alive
}
func FullHealthFilter(game Game, actor Actor, context Context) bool {
	return actor.Damage == 0
}
func ActiveFilter(game Game, actor Actor, context Context) bool {
	return actor.IsActive()
}
func InactiveFilter(game Game, actor Actor, context Context) bool {
	return actor.PositionID == nil
}
func PlayerFilter(game Game, actor Actor, context Context) bool {
	if context.SourcePlayerID == nil {
		return false
	}
	return actor.PlayerID == *context.SourcePlayerID
}
func SourceFilter(game Game, actor Actor, context Context) bool {
	if context.SourceActorID == nil {
		return false
	}
	return actor.ID == *context.SourceActorID
}
func ParentFilter(game Game, actor Actor, context Context) bool {
	if context.ParentActorID == nil {
		return false
	}
	return actor.ID == *context.ParentActorID
}
func TargetableFilter(game Game, actor Actor, context Context) bool {
	return ComposeAF(
		AliveFilter,
		ActiveFilter,
		func(game Game, actor Actor, context Context) bool {
			/*
				resolved := actor.Resolve(game)

				if resolved.State == ActorStateIncorporeal {
					return false
				}
			*/
			return true
		},
	)(game, actor, context)
}
func TargetFilter(game Game, actor Actor, context Context) bool {
	if slices.Contains(context.TargetActorIDs, actor.ID) {
		return true
	}
	if actor.IsActive() && slices.Contains(context.TargetPositionIDs, *actor.PositionID) {
		return true
	}

	return false
}
func ActiveTransactionFilter(game Game, actor Actor, context Context) bool {
	if game.ActiveTransaction == nil {
		return false
	}

	ctx := game.ActiveTransaction.Context
	if ctx.SourceActorID == nil || context.SourceActorID == nil {
		return false
	}
	context_source := *ctx.SourceActorID == *context.SourceActorID
	context_target := false
	actor_source := *ctx.SourceActorID == actor.ID
	actor_target := false

	targets := game.GetTargets(ctx)
	for _, t := range targets {
		if t.ID == *context.SourceActorID {
			context_target = true
		}
		if t.ID == actor.ID {
			actor_target = t.ID == *context.SourceActorID
		}
	}

	return (context_source || context_target) && (actor_source || actor_target)
}
func ActiveTargetFilter(game Game, actor Actor, context Context) bool {
	if game.ActiveTransaction == nil {
		return false
	}

	atx_ctx := game.ActiveTransaction.Context
	if atx_ctx.SourceActorID == nil || context.SourceActorID == nil {
		return false
	}
	context_target := false
	actor_target := false

	targets := game.GetTargets(atx_ctx)
	for _, t := range targets {
		if t.ID == *context.SourceActorID {
			context_target = true
		}
		if t.ID == actor.ID {
			actor_target = t.ID == *context.SourceActorID
		}
	}

	return context_target && actor_target
}
func ActiveTargetTeamSourceFilter(game Game, actor Actor, context Context) bool {
	if game.ActiveTransaction == nil {
		return false
	}

	ctx := game.ActiveTransaction.Context
	if ctx.SourcePlayerID == nil || context.SourcePlayerID == nil {
		return false
	}

	actor_player := *ctx.SourcePlayerID == actor.PlayerID
	context_player := *ctx.SourcePlayerID == *context.SourcePlayerID
	return actor_player || context_player
}
func TeamFilter(game Game, actor Actor, context Context) bool {
	if context.SourcePlayerID == nil {
		return false
	}
	return actor.PlayerID == *context.SourcePlayerID
}
func OtherTeamFilter(game Game, actor Actor, context Context) bool {
	if context.SourcePlayerID == nil {
		return false
	}
	return actor.PlayerID != *context.SourcePlayerID
}
func IsAtOrBelowHealthRatio(ratio float64) func(Game, Actor, Context) bool {
	return func(game Game, actor Actor, context Context) bool {
		hp := float64(actor.Stats[StatHP])
		damage := float64(actor.Damage)
		return ratio >= (hp-damage)/hp
	}
}
func HasAppliedModifier(modifierID uuid.UUID) func(Game, Actor, Context) bool {
	return func(game Game, actor Actor, context Context) bool {
		for mid, _ := range actor.AppliedModifiers {
			if mid == modifierID {
				return true
			}
		}

		return false
	}
}
func GameHasWeather(weather GameWeather) func(Game, Actor, Context) bool {
	return func(g Game, a Actor, ctx Context) bool {
		state, _ := g.GetState(ctx)
		return state.Weather == weather
	}
}
func GameHasTerrain(terrain GameTerrain) func(Game, Actor, Context) bool {
	return func(g Game, a Actor, ctx Context) bool {
		state, _ := g.GetState(ctx)
		return state.Terrain == terrain
	}
}
func HasState(state ActorStateType) func(Game, Actor, Context) bool {
	return func(g Game, a Actor, ctx Context) bool {
		return a.State == state
	}
}

/**
 * RESOLVED FILTERS
 * These filters required modifiers to be resolved to check things like
 * protected, chakra amount, and things that can change with modifers
 * THESE CANNOT BE USED IN MODIFIERS
 */
func RHasChakraFilter(amount int) func(Game, Actor, Context) bool {
	return func(game Game, actor Actor, context Context) bool {
		resolved := actor.Resolve(game)
		return resolved.HasChakra(amount)
	}
}
func RIsProtectedFilter(protected bool) func(Game, Actor, Context) bool {
	return func(game Game, actor Actor, context Context) bool {
		resolved := actor.Resolve(game)
		return resolved.Protected == protected
	}
}
func RHasState(state ActorStateType) func(Game, Actor, Context) bool {
	return func(g Game, a Actor, ctx Context) bool {
		resolved := a.Resolve(g)
		return resolved.State == state
	}
}
