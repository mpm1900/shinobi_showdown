package game

import (
	"slices"

	"github.com/google/uuid"
)

type Context struct {
	ActionID          *uuid.UUID     `json:"action_ID"`
	ModifierID        *uuid.UUID     `json:"modifier_ID"`
	SourcePlayerID    *uuid.UUID     `json:"source_player_ID"`
	ParentActorID     *uuid.UUID     `json:"parent_actor_ID"`
	SourceActorID     *uuid.UUID     `json:"source_actor_ID"`
	TargetActorIDs    []uuid.UUID    `json:"target_actor_IDs"`
	TargetPositionIDs []uuid.UUID    `json:"target_position_IDs"`
	Meta              map[string]int `json:"-"`
}

func NewContext() Context {
	return Context{
		TargetActorIDs:    []uuid.UUID{},
		TargetPositionIDs: []uuid.UUID{},
	}
}
func MakeContextForActor(actor Actor) Context {
	return Context{
		SourcePlayerID:    &actor.PlayerID,
		SourceActorID:     &actor.ID,
		ParentActorID:     &actor.ID,
		TargetActorIDs:    []uuid.UUID{actor.ID},
		TargetPositionIDs: []uuid.UUID{},
	}
}
func (c *Context) SetMeta(key string, value int) {
	if c.Meta == nil {
		c.Meta = map[string]int{}
	}

	c.Meta[key] = value
}
func (c *Context) GetMeta(key string) (int, bool) {
	if c.Meta == nil {
		c.Meta = map[string]int{}
	}

	value, ok := c.Meta[key]
	return value, ok
}

func ResolveModifierTransactionContext(
	fallback Context,
	transactions []Transaction[Modifier],
	transactionID *uuid.UUID,
) Context {
	if transactionID == nil {
		return fallback
	}

	for _, transaction := range transactions {
		if transaction.ID == *transactionID {
			return transaction.Context
		}
	}

	return fallback
}
func (c Context) WithPlayer(playerID uuid.UUID) Context {
	c.SourcePlayerID = &playerID
	return c
}
func (c Context) WithSource(sourceID uuid.UUID) Context {
	c.SourceActorID = &sourceID
	return c
}
func (c Context) WithTargetIDs(targetActorIDs []uuid.UUID) Context {
	c.TargetActorIDs = targetActorIDs
	return c
}

func (c *Context) FilterOutTarget(actor Actor) {
	c.TargetActorIDs = slices.DeleteFunc(c.TargetActorIDs, func(ID uuid.UUID) bool {
		return ID == actor.ID
	})
	if actor.PositionID != nil {
		c.TargetPositionIDs = slices.DeleteFunc(c.TargetPositionIDs, func(posID uuid.UUID) bool {
			return posID == *actor.PositionID
		})
	}
}

func (g Game) GetTargets(context Context) []Actor {
	targets := g.GetActors(func(a Actor) bool {
		return slices.Contains(context.TargetActorIDs, a.ID)
	})
	posTargets := g.GetActors(func(a Actor) bool {
		if a.PositionID == nil {
			return false
		}

		return slices.Contains(context.TargetPositionIDs, *a.PositionID)
	})

	targets = append(targets, posTargets...)
	return targets
}

func (g Game) GetSource(context Context) (Actor, bool) {
	if context.SourceActorID == nil {
		return Actor{}, false
	}

	return g.GetActorByID(*context.SourceActorID)
}

func (g Game) GetParent(context Context) (Actor, bool) {
	if context.ParentActorID == nil {
		return Actor{}, false
	}

	return g.GetActorByID(*context.ParentActorID)
}
