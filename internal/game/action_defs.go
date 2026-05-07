package game

import (
	"math/rand/v2"

	"github.com/google/uuid"
)

func NotNewlyInactiveFilter(game Game, actor Actor, context Context) bool {
	return actor.InactiveTurns > 0
}

func SwitchFilter(game Game, actor Actor, context Context) bool {
	possible := game.GetActorsFilters(context, ComposeAF(TeamFilter, InactiveFilter, AliveFilter, NotNewlyInactiveFilter))
	if len(possible) == 0 {
		return true
	}

	return NotNewlyInactiveFilter(game, actor, context)
}

func NewNoopPlayer(groupID *uuid.UUID) ActorMutation {
	return MakeActorMutation(
		groupID,
		0,
		PlayerFilter,
		func(g Game, a Actor, c Context) Actor {
			return a
		},
	)
}
func NewNoopSource(groupID *uuid.UUID) ActorMutation {
	return MakeActorMutation(
		groupID,
		0,
		SourceFilter,
		func(g Game, a Actor, c Context) Actor {
			return a
		},
	)
}
func NewNoopParent(groupID *uuid.UUID) ActorMutation {
	return MakeActorMutation(
		groupID,
		0,
		ParentFilter,
		func(g Game, a Actor, c Context) Actor {
			return a
		},
	)
}
func NewNoopActive(groupID *uuid.UUID) ActorMutation {
	return MakeActorMutation(
		groupID,
		0,
		ActiveFilter,
		func(g Game, a Actor, c Context) Actor {
			return a
		},
	)
}

var SwitchPositions = GameMutation{
	Delta: func(p Game, g Game, context Context) Game {
		source, ok := g.GetSource(context)
		if !ok || source.PositionID == nil {
			return g
		}

		targets := g.GetTargets(context)
		if len(targets) == 0 {
			return g
		}

		g.SetPosition(targets[0], source.PositionID)
		return g
	},
}

var SetPositions = GameMutation{
	Delta: func(p Game, g Game, context Context) Game {
		targets := g.GetTargets(context)
		positionCount := len(context.TargetPositionIDs)
		if len(targets) == 0 || positionCount == 0 {
			return g
		}

		limit := min(len(targets), positionCount)
		for i := range limit {
			positionID := context.TargetPositionIDs[i]
			g.SetPosition(targets[i], &positionID)
		}

		return g
	},
}

var RemovePositions = GameMutation{
	Delta: func(p Game, g Game, context Context) Game {
		targets := g.GetTargets(context)
		for i := range len(targets) {
			g.SetPosition(targets[i], nil)
		}

		return g
	},
}

var Switch = Action{
	ID: uuid.MustParse("2c6bf2e3-4049-4f1f-b18a-f4fd844ef31a"),
	Config: ActionConfig{
		Name:        "Switch",
		Description: "Switches user out and target ally into battle.",
	},
	Meta: ActionMeta{
		Switch: true,
	},
	Locked:          true,
	TargetType:      TargetActorID,
	TargetPredicate: ComposeAF(TeamFilter, InactiveFilter, AliveFilter, SwitchFilter),
	ContextValidate: TargetLengthFilter(1),
	ActionMutation: ActionMutation{
		Priority: ActionPrioritySwitch,
		Filter:   TrueGameFilter,
		Delta: func(p Game, input Game, context Context) []Transaction[GameMutation] {
			transactions := []GameTransaction{
				MakeTransaction(SwitchPositions, context),
			}
			return transactions
		},
	},
}

var SwitchInIds map[int]uuid.UUID = map[int]uuid.UUID{
	1: uuid.MustParse("f4e25f49-35c1-405d-a34a-c8ba22f595f8"),
	2: uuid.MustParse("11a552fc-a768-4f58-aef1-fecba66d95f8"),
}

func SwitchIn(count int) Action {
	return Action{
		ID: SwitchInIds[count],
		Config: ActionConfig{
			Name: "Switch In",
		},
		Meta: ActionMeta{
			Switch: true,
		},
		Locked:          true,
		TargetType:      TargetActorID,
		TargetPredicate: ComposeAF(TeamFilter, InactiveFilter, AliveFilter, SwitchFilter),
		ContextValidate: TargetLengthFilter(count),
		ActionMutation: ActionMutation{
			Priority: ActionPrioritySwitch,
			Filter:   TrueGameFilter,
			Delta: func(parent Game, input Game, context Context) []Transaction[GameMutation] {
				transactions := []GameTransaction{}
				min := min(len(context.TargetActorIDs), len(context.TargetPositionIDs))
				for i := range min {
					c := context
					c.TargetActorIDs = []uuid.UUID{c.TargetActorIDs[i]}
					c.TargetPositionIDs = []uuid.UUID{c.TargetPositionIDs[i]}
					transactions = append(transactions, MakeTransaction(SetPositions, c))
				}

				return transactions
			},
		},
	}
}

func makeAttack(
	ID uuid.UUID,
	config ActionConfig,
	with func(g Game, context Context, transactions []GameTransaction) []GameTransaction,
) Action {
	return Action{
		ID:              ID,
		Config:          config,
		TargetType:      TargetPositionID,
		TargetPredicate: ComposeAF(OtherFilter, TargetableFilter),
		ContextValidate: PositionsLengthFilter(*config.TargetCount),
		ActionMutation: ActionMutation{
			Priority: ActionPriorityDefault,
			Filter: ComposeGF(
				SourceIsAlive,
			),
			Delta: func(p Game, g Game, context Context) []GameTransaction {
				transactions := []GameTransaction{}

				conf, _ := GetActiveActionConfig(g, config)
				crit_result := MakeCriticalCheck(conf)
				damages := NewDamage(conf, NewDamageConfig(crit_result.Ratio, RandomDamageFactor()))
				transactions = append(
					transactions,
					MakeDamageTransactions(context, damages)...,
				)

				return with(g, context, transactions)
			},
		},
	}
}

var struggleID = uuid.MustParse("33ac9155-7427-4774-bc32-2d3161fa9b47")
var Struggle = MakeStruggle()

func MakeStruggle() Action {
	config := ActionConfig{
		Name:        "Struggle",
		Description: "Deals 1/4th HP in recoil damage. Can be used when no other actions are available.",
		TargetCount: Ptr(0),
		Nature:      Ptr(NsPure),
		Power:       Ptr(50),
		Stat:        Ptr(StatAttack),
		Jutsu:       Taijutsu,
		Cooldown:    Ptr(0),
		CritChance:  Ptr(5),
		CritMod:     1.5,
	}

	action := makeAttack(
		struggleID,
		config,
		func(g Game, context Context, transactions []GameTransaction) []GameTransaction {
			source, ok := g.GetSource(context)
			if !ok {
				return transactions
			}

			recoilDamage := RatioDamage(0.25)
			recoilContext := MakeContextForActor(source)
			transactions = append(
				transactions,
				MakeDamageTransactions(recoilContext, recoilDamage)...,
			)

			return transactions
		},
	)
	action.Meta.Struggle = true
	action.TargetPredicate = NoneFilter
	action.MapContext = func(g Game, context Context) Context {
		other_actors := g.GetActorsFilters(context, ComposeAF(ActiveFilter, OtherTeamFilter))
		if len(other_actors) == 0 {
			return context
		}

		rand_index := rand.IntN(len(other_actors))
		context.TargetPositionIDs = append(context.TargetPositionIDs, *other_actors[rand_index].PositionID)
		return context
	}

	return action
}

var cancelSummonID = uuid.MustParse("a0491437-2cda-407e-9fc8-4022b2e1f5e3")
var CancelSummon = Action{
	ID: cancelSummonID,
	Config: ActionConfig{
		Name:        "Cancel Summon",
		Description: "Unsummons the user's summon.",
	},
	TargetType: TargetActorID,
	TargetPredicate: func(g Game, a Actor, ctx Context) bool {
		return false
	},
	ContextValidate: func(ctx Context) bool {
		return true
	},
	ActionMutation: ActionMutation{
		Priority: ActionPriorityDefault,
		Filter:   SourceIsAlive,
		Delta: func(p, g Game, context Context) []Transaction[GameMutation] {
			transactions := []GameTransaction{}

			mut := GameMutation{
				Delta: func(mp, mg Game, mc Context) Game {
					mg.UpdateActor(*mc.SourceActorID, func(a Actor) Actor {
						a.SetSummon(nil)
						return a
					})
					return mg
				},
			}

			transactions = append(
				transactions,
				MakeTransaction(mut, context),
			)

			return transactions
		},
	},
}
