package game

import (
	"slices"

	"github.com/google/uuid"
)

/**
 * Context Filters
 */

type ContextFilter func(Context) bool

func ComposeCF(filters ...ContextFilter) ContextFilter {
	return func(context Context) bool {
		for _, filter := range filters {
			if !filter(context) {
				return false
			}
		}

		return true
	}
}

func TargetLengthFilter(length int) func(Context) bool {
	return func(context Context) bool {
		return len(context.TargetActorIDs) == length
	}
}
func PositionsLengthFilter(length int) func(Context) bool {
	return func(context Context) bool {
		return len(context.TargetPositionIDs) == length
	}
}

/**
 * Game Filters
 */

type GameFilter func(Game, Game, Context) bool

func ComposeGF(filters ...GameFilter) GameFilter {
	return func(parent Game, game Game, context Context) bool {
		for _, filter := range filters {
			if !filter(parent, game, context) {
				return false
			}
		}

		return true
	}
}

func TrueGameFilter(parent Game, game Game, context Context) bool {
	return true
}
func FalseGameFilter(parent Game, game Game, context Context) bool {
	return false
}
func HasWeather(weather GameWeather) func(Game, Game, Context) bool {
	return func(parent, game Game, context Context) bool {
		state, _ := game.GetState(context)
		return state.Weather == weather
	}
}
func SourceIsAlive(parent Game, game Game, context Context) bool {
	source, ok := game.GetSource(context)
	if !ok {
		return false
	}

	return source.Alive
}
func SourceIsActionOffCooldown(parent Game, game Game, context Context) bool {
	if context.ActionID == nil {
		return false
	}

	source, ok := game.GetSource(context)
	if !ok {
		return false
	}

	action, ok := source.GetActionByID(game, *context.ActionID)
	if !ok {
		return true
	}

	return action.Cooldown == nil
}
func SourceHasActiveTurns(turns int) func(Game, Game, Context) bool {
	return func(parent Game, g Game, context Context) bool {
		if context.ActionID == nil {
			return false
		}

		source, ok := g.GetSource(context)
		if !ok {
			return false
		}

		return source.ActiveTurns == turns
	}
}
func SourceIsAtOrBelowHealth(ratio float64) func(Game, Game, Context) bool {
	return func(parent, g Game, context Context) bool {
		source, ok := g.GetSource(context)
		if !ok {
			return false
		}
		resolved := source.Resolve(g)
		hp := float64(resolved.Stats[StatHP])
		damage := float64(resolved.Damage)
		return ratio >= (hp-damage)/hp
	}
}
func SourceHasHpRatio(ratio float64) func(Game, Game, Context) bool {
	return func(p, g Game, context Context) bool {
		s, ok := g.GetSource(context)
		if !ok {
			return false
		}
		source := s.Resolve(g)
		hp := source.Stats[StatHP]
		return float64(hp-source.Damage)/float64(hp) > ratio
	}
}

func TargetsAreActive(parent Game, game Game, context Context) bool {
	targets := game.GetTargets(context)
	for _, target := range targets {
		if target.PositionID != nil {
			return true
		}
	}
	return false
}
func TargetsAreOneAlive(parent Game, game Game, context Context) bool {
	targets := game.GetTargets(context)
	for _, target := range targets {
		if target.Alive {
			return true
		}
	}
	return false
}
func TargetsHaveAppliedModifier(modifierID uuid.UUID) func(Game, Game, Context) bool {
	return func(parent Game, game Game, context Context) bool {
		targets := game.GetTargets(context)
		for _, target := range targets {
			resolved := target.Resolve(game)
			for mid, _ := range resolved.AppliedModifiers {
				if mid == modifierID {
					return true
				}
			}
		}
		return false
	}
}

/**
 * Trigger filters
 */

type TriggerFilter = func(Game, Game, Context, Transaction[Modifier]) bool

func ComposeTF(filters ...TriggerFilter) TriggerFilter {
	return func(parent Game, game Game, context Context, tx Transaction[Modifier]) bool {
		for _, filter := range filters {
			if !filter(parent, game, context, tx) {
				return false
			}
		}

		return true
	}
}
func Match__Player_Player(p, g Game, context Context, tx Transaction[Modifier]) bool {
	if context.SourcePlayerID == nil || tx.Context.SourcePlayerID == nil {
		return false
	}

	return *context.SourcePlayerID == *tx.Context.SourcePlayerID
}
func Match__TargetActor_SourceActor(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
	targets := game.GetTargets(context)
	if len(targets) == 0 || modifier_tx.Context.SourceActorID == nil {
		return false
	}

	for _, t := range targets {
		if t.ID == *modifier_tx.Context.SourceActorID {
			return true
		}
	}
	return false
}
func Match__SourceActor_TargetActor(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
	targets := game.GetTargets(modifier_tx.Context)
	if len(targets) == 0 || context.SourceActorID == nil {
		return false
	}

	for _, t := range targets {
		if t.ID == *context.SourceActorID {
			return true
		}
	}
	return false
}
func Match__SourceActor_SourceActor(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
	if context.SourceActorID == nil || modifier_tx.Context.SourceActorID == nil {
		return false
	}

	return *context.SourceActorID == *modifier_tx.Context.SourceActorID
}
func NotMatch__SourceActor_SourceActor(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
	if context.SourceActorID == nil || modifier_tx.Context.SourceActorID == nil {
		return false
	}

	return *context.SourceActorID != *modifier_tx.Context.SourceActorID
}
func Source__IsAtOrBelowHealth(ratio float64) TriggerFilter {
	return func(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
		return SourceIsAtOrBelowHealth(ratio)(parent, game, modifier_tx.Context)
	}
}
func Source__IsActive(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
	source, ok := game.GetSource(modifier_tx.Context)
	if !ok {
		return false
	}

	return source.IsActive()
}
func Modifier__IsOneOf(modifierIDs ...uuid.UUID) TriggerFilter {
	return func(parent Game, game Game, context Context, modifier_tx Transaction[Modifier]) bool {
		if context.ModifierID == nil {
			return false
		}

		modifier, ok := game.GetModifierTxByID(*context.ModifierID)
		if !ok || modifier.Mutation.GroupID == nil {
			return false
		}

		return slices.Contains(modifierIDs, *modifier.Mutation.GroupID)
	}
}
