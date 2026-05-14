package game

import (
	"fmt"
	"math/rand/v2"

	"github.com/google/uuid"
)

const (
	ActionPrioritySwitch  = 10
	ActionPriorityP5      = 5
	ActionPriorityProtect = 4
	ActionPriorityP3      = 3
	ActionPriorityP2      = 2
	ActionPriorityP1      = 1
	ActionPriorityDefault = 0
	ActionPrioritySlow    = -1
	ActionPrioritySlow2   = -2
	ActionPrioritySlow3   = -3
)

type ActionJutsu string

const (
	Bukijutsu ActionJutsu = "bukijutsu"
	Fuinjutsu ActionJutsu = "fuinjutsu"
	Genjutsu  ActionJutsu = "genjutsu"
	Ninjutsu  ActionJutsu = "ninjutsu"
	Senjutsu  ActionJutsu = "senjutsu"
	Taijutsu  ActionJutsu = "taijutsu"
)

type ActionConfig struct {
	Accuracy       *int        `json:"accuracy,omitempty"`
	Cooldown       *int        `json:"cooldown,omitempty"`
	Cost           *int        `json:"cost,omitempty"`
	CritChance     *int        `json:"-"`
	CritMod        float64     `json:"-"`
	LifeSteal      *float64    `json:"-"`
	Name           string      `json:"name"`
	Nature         *NatureSet  `json:"nature,omitempty"`
	Power          *int        `json:"power,omitempty"`
	Recoil         *float64    `json:"-"`
	Stat           *ActorStat  `json:"stat,omitempty"`
	TargetCount    *int        `json:"-"`
	Jutsu          ActionJutsu `json:"jutsu"`
	Description    string      `json:"description"`
	LogSuccess     *string     `json:"log_success"`
	LogFailure     *string     `json:"-"`
	IgnoreRedirect bool        `json:"-"`
	SubPriority    int         `json:"-"`
}

type ActionTargetType string

const (
	TargetActorID    ActionTargetType = "target-actor-id"
	TargetPositionID ActionTargetType = "target-position-type"
)

type ActionMutation Mutation[Game, Game, []Transaction[GameMutation]]

type ActionMeta struct {
	Switch   bool `json:"switch"`
	Struggle bool `json:"struggle"`
}

/** [This comment was not written by an LLM]
 * Action Function Members for Action "a"
 *
 * action.Filter(Game, *Context) => can this action be taken with this context.
 * -- this is often done for a chakra or disabled check
 *
 * action.TargetPredicate(Actor, *Context) => is this actor a valid target for this action
 * -- this is effectively the "targets generator" for users to choose.
 *
 * action.ContextValidate(*Context) => does this context represent a valid targets selection for this action
 * -- this is used to check "is the number of targets correct?" and other checks.
 *
 *
 * action.Delta(Game, *Context) => resolution of the Action
 * -- can include random events
 */
type Action struct {
	ActionMutation
	ID              uuid.UUID                       `json:"ID"`
	Config          ActionConfig                    `json:"config"`
	Summon          bool                            `json:"summon"`
	Locked          bool                            `json:"locked"`
	Disabled        bool                            `json:"disabled"`
	TargetType      ActionTargetType                `json:"target_type"`
	TargetPredicate func(Game, Actor, Context) bool `json:"-"`
	ContextValidate func(Context) bool              `json:"-"`
	MapContext      func(Game, Context) Context     `json:"-"`
	Cost            GameMutation                    `json:"-"`
	Cooldown        *int                            `json:"cooldown"`
	Meta            ActionMeta                      `json:"meta"`
}

func ResolveAction(game *Game, transaction Transaction[Action]) []GameTransaction {
	action := transaction.Mutation
	context := transaction.Context

	if context.SourceActorID != nil {
		if queue, ok := game.QueuedActions[*context.SourceActorID]; ok {
			delete(game.QueuedActions, *context.SourceActorID)
			if queue.Mutation != transaction.Mutation.ID {
				fmt.Println("ERROR: INVALID ACTION EXECUTED")
				return []GameTransaction{}
			}
		}
	}

	source, ok := game.GetSource(context)
	if !ok && context.SourceActorID != nil {
		return []GameTransaction{}
	}

	/**
	 * Source Can-Act Checks
	 */
	if ok && !action.Meta.Switch {
		resolved := source.Resolve(*game)
		if !resolved.CanAct(game, context) {
			return []GameTransaction{}
		}
	}

	if transaction.Context.SourceActorID != nil {
		cost := transaction.Mutation.Cost
		if cost.Delta != nil {
			costTx := MakeTransaction(cost, transaction.Context)
			game.Transactions = append(game.Transactions, costTx)
		}
	}

	/**
	 * Action Can-Act Checks
	 */
	if action.Disabled || !action.Filter(*game, *game, transaction.Context) {
		context.ActionID = &action.ID
		logStart := NewLogContext("$source$ used $action$", context)
		logFail := NewLogContext("$action$ failed.", context)
		if action.Config.LogFailure != nil {
			logFail = NewLogContext(*action.Config.LogFailure, context)
		}

		game.PushLog(logStart)
		game.PushLog(logFail)
		return []GameTransaction{}
	}

	if action.Config.Cooldown != nil && *action.Config.Cooldown > 0 {
		game.SetActionCooldown(
			*transaction.Context.SourceActorID,
			action.ID,
			*action.Config.Cooldown,
		)
	}

	if action.MapContext != nil {
		context = action.MapContext(*game, context)
	}

	if ok {
		if !action.Meta.Struggle {
			game.UpdateActor(source.ID, func(a Actor) Actor {
				a.LastUsedActionTX = &transaction
				return a
			})
		}
		log := NewLogContext("$source$ used $action$", context)
		if action.Config.LogSuccess != nil {
			log = NewLogContext(*action.Config.LogSuccess, context)
		}
		game.PushLog(log)
	}

	return action.Delta(*game, *game, context)
}

func GetAccuracy(game Game, source ResolvedActor, target ResolvedActor, ignoreModifiers bool) float64 {
	ratio := float64(source.Stats[StatAccuracy]) / float64(target.Stats[StatEvasion])
	fmt.Printf("acc ratio = %d / %d = %f \n", source.Stats[StatAccuracy], target.Stats[StatEvasion], ratio)
	return ratio
}

func MakeActionRoll() int {
	return rand.IntN(100)
}

type ChanceResult struct {
	Chance  int
	Roll    int
	Success bool
	Ratio   float64
}

func MakeCriticalCheck(action ActionConfig) ChanceResult {
	if action.CritChance == nil {
		return ChanceResult{
			Success: false,
		}
	}

	fmt.Println("CRIT CHANCE:", *action.CritChance)
	threshold := int(*action.CritChance)
	roll := MakeActionRoll()
	success := roll <= threshold
	ratio := 1.0
	if success {
		ratio = action.CritMod
	}

	return ChanceResult{
		Chance:  threshold,
		Roll:    roll,
		Success: success,
		Ratio:   ratio,
	}
}

func MakeAccuracyCheck(g *Game, action ActionConfig, source ResolvedActor, target ResolvedActor, ignoreModifiers bool) ChanceResult {
	base_accuracy := GetAccuracy(*g, source, target, ignoreModifiers)
	if action.Accuracy == nil {
		return ChanceResult{
			Success: true,
		}
	}

	accuracy := Round(base_accuracy * (float64(*action.Accuracy + source.ActionAccuracyOffset)))
	roll := MakeActionRoll()
	fmt.Printf("acc = %d, roll = %d\n", accuracy, roll)
	return ChanceResult{
		Chance:  accuracy,
		Roll:    roll,
		Success: roll <= accuracy,
	}
}

func GetActiveActionConfig(g Game, fallback ActionConfig) (ActionConfig, bool) {
	if g.ActiveTransaction == nil {
		return fallback, false
	}

	return g.ActiveTransaction.Mutation.Config, true
}
