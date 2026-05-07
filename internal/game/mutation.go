package game

import (
	"slices"

	"github.com/google/uuid"
)

type Mutation[P any, I any, O any] struct {
	Delta    func(parent P, input I, context Context) O    `json:"-"`
	Filter   func(parent P, input I, context Context) bool `json:"-"`
	Priority int                                           `json:"priority"`
}

/**
 * Game Mutations
 * [GameMutation]
 */
type GameMutation = Mutation[Game, Game, Game]

func NewGameMutation() GameMutation {
	return GameMutation{
		Delta:    func(parent Game, input Game, context Context) Game { return input },
		Filter:   func(parent Game, input Game, context Context) bool { return true },
		Priority: 0,
	}
}

func ComposeGameMutations(mutations ...GameMutation) GameMutation {
	if len(mutations) == 0 {
		return NewGameMutation()
	}

	slices.SortFunc(mutations, func(a, b GameMutation) int {
		return b.Priority - a.Priority
	})

	return GameMutation{
		Delta: func(g Game, input Game, context Context) Game {
			for _, mut := range mutations {
				if mut.Delta != nil {
					input = mut.Delta(g, input, context)
				}
			}

			return input
		},
		Filter: func(g Game, input Game, context Context) bool {
			for _, mut := range mutations {
				if mut.Filter != nil && !mut.Filter(g, input, context) {
					return false
				}
			}
			return true
		},
		Priority: mutations[0].Priority,
	}
}

func AddLogs(logs ...GameLog) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			g.Log = append(g.Log, logs...)

			return g
		},
	}
}

/**
 * Actor Mutations
 * [ActorMutation]
 */
type ActorMutation struct {
	Mutation[Game, Actor, Actor]
	ModifierGroupID *uuid.UUID
	TransactionID   *uuid.UUID
}

func MakeActorMutation(
	modifierGroupID *uuid.UUID,
	priority int,
	filter func(Game, Actor, Context) bool,
	delta func(Game, Actor, Context) Actor,
) ActorMutation {
	return ActorMutation{
		ModifierGroupID: modifierGroupID,
		Mutation: Mutation[Game, Actor, Actor]{
			Filter:   filter,
			Delta:    delta,
			Priority: priority,
		},
	}
}

/**
 * GameState Mutations
 * [GameStateMutation]
 */
type GameStateMutation struct {
	Mutation[Game, GameState, GameState]
	ModifierGroupID *uuid.UUID
	TransactionID   *uuid.UUID
}

func MakeGameStateMutation(
	modifierGroupID *uuid.UUID,
	priority int,
	filter func(Game, GameState, Context) bool,
	delta func(Game, GameState, Context) GameState,
) GameStateMutation {
	return GameStateMutation{
		ModifierGroupID: modifierGroupID,
		Mutation: Mutation[Game, GameState, GameState]{
			Filter:   filter,
			Delta:    delta,
			Priority: priority,
		},
	}
}
