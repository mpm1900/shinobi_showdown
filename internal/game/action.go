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

type ActionTargetType string

const (
	TargetActorID    ActionTargetType = "target-actor-id"
	TargetPositionID ActionTargetType = "target-position-type"
)

type ActionConfig struct {
	Accuracy       *int             `json:"accuracy,omitempty"`
	Cooldown       *int             `json:"cooldown,omitempty"`
	Cost           *int             `json:"cost,omitempty"`
	CritChance     *int             `json:"-"`
	CritStage      *int             `json:"-"`
	CritMod        float64          `json:"-"`
	LifeSteal      *float64         `json:"-"`
	Name           string           `json:"name"`
	Nature         *NatureSet       `json:"nature,omitempty"`
	Power          *int             `json:"power,omitempty"`
	Recoil         *float64         `json:"-"`
	Stat           *ActorStat       `json:"stat,omitempty"`
	TargetCount    *int             `json:"-"`
	TargetType     ActionTargetType `json:"target_type"`
	Jutsu          ActionJutsu      `json:"jutsu"`
	Description    string           `json:"description"`
	LogSuccess     *string          `json:"log_success"`
	LogFailure     *string          `json:"log_failure"`
	IgnoreRedirect bool             `json:"-"`
	SubPriority    int              `json:"-"`
	Summon         bool             `json:"summon"`
	Switch         bool             `json:"switch"`
	Struggle       bool             `json:"-"`
}

func (ac ActionConfig) With(f func(*ActionConfig)) ActionConfig {
	f(&ac)
	return ac
}

type ActionMutation Mutation[Game, Game, []Transaction[GameMutation]]

type ActionState struct {
	Disabled bool `json:"disabled"`
	Cooldown *int `json:"cooldown"`
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
	ID     uuid.UUID    `json:"ID"`
	Config ActionConfig `json:"config"`
	State  ActionState  `json:"state"`

	TargetPredicate func(Game, Actor, Context) bool `json:"-"`
	ContextValidate func(Context) bool              `json:"-"`
	MapContext      func(Game, Context) Context     `json:"-"`
	Cost            GameMutation                    `json:"-"`
}

func ResolveAction(game *Game, transaction Transaction[Action]) []GameTransaction {
	action := transaction.Mutation
	context := transaction.Context
	transactions := NewTransactionBuilder()

	if context.SourceActorID != nil {
		if queue, ok := game.QueuedActions[*context.SourceActorID]; ok {
			delete(game.QueuedActions, *context.SourceActorID)
			if queue.Mutation != transaction.Mutation.ID {
				fmt.Println("ERROR: INVALID ACTION EXECUTED")
				return transactions.Build()
			}
		}
	}

	source, hasSource := game.GetSource(context)
	if !hasSource && context.SourceActorID != nil {
		return transactions.Build()
	}

	/**
	 * Source Can-Act Checks
	 */
	if hasSource && !action.Config.Switch {
		resolved := source.Resolve(*game)
		if !resolved.CanAct(game, context) {
			return transactions.Build()
		}
	}

	cost := action.Cost
	if cost.Delta != nil {
		costTx := MakeTransaction(cost, transaction.Context)
		transactions.PushOne(costTx)
	}

	context.ActionID = &action.ID
	log := NewLogContext("$source$ used $action$", context)
	if action.Config.LogSuccess != nil {
		log = NewLogContext(*action.Config.LogSuccess, context)
	}

	/**
	 * Action Can-Act Checks
	 */
	if action.State.Disabled || !action.Filter(*game, *game, transaction.Context) {
		logFail := NewLogContext("$action$ failed.", context)
		if action.Config.LogFailure != nil {
			logFail = NewLogContext(*action.Config.LogFailure, context)
		}

		if action.Config.LogSuccess != nil {
			log = NewLogContext(*action.Config.LogSuccess, context)
		}

		transactions.PushOne(MakeTransaction(AddLogs(log), context))
		transactions.PushOne(MakeTransaction(AddLogs(logFail), context))
		return transactions.Build()
	}

	if hasSource {
		transactions.PushOne(MakeTransaction(AddLogs(log), context))
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

	if !action.Config.Struggle {
		game.UpdateActor(source.ID, func(a Actor) Actor {
			a.LastUsedActionTX = &transaction
			return a
		})
	}

	transactions.Push(action.Delta(*game, *game, context))
	return transactions.Build()
}

func GetAccuracy(game Game, source ResolvedActor, target ResolvedActor, ignoreModifiers bool) float64 {
	if ignoreModifiers {
		ratio := float64(source.UnmodifiedStats[StatAccuracy]) / float64(target.UnmodifiedStats[StatEvasion])
		return ratio
	}

	ratio := float64(source.Stats[StatAccuracy]) / float64(target.Stats[StatEvasion])
	return ratio
}

func GetCriticalChance(action ActionConfig, source ResolvedActor) float64 {
	crit_chance := 0.0
	if action.CritChance == nil && action.CritStage == nil {
		return crit_chance
	}

	if action.CritChance != nil {
		crit_chance = float64(*action.CritChance)
	}
	if action.CritStage != nil {
		stage := *action.CritStage + source.Stages[StatCritical]
		crit_chance = getCriticalStage(stage)
	}

	crit_chance += float64(source.BaseStats[StatCritical])

	return crit_chance
}

func MakeActionRoll() int {
	return rand.IntN(100)
}

type ChanceResult struct {
	Chance  float64
	Roll    int
	Success bool
	Ratio   float64
}

var criticalStages = map[int]float64{
	0: CRITICAL_STAGE_0,
	1: CRITICAL_STAGE_1,
	2: CRITICAL_STAGE_2,
	3: CRITICAL_STAGE_3,
}

func getCriticalStage(stage int) float64 {
	stage = max(0, stage)
	stage = min(stage, len(criticalStages)-1)
	return criticalStages[stage]
}

func MakeCriticalCheck(action ActionConfig, source ResolvedActor) ChanceResult {
	crit_chance := GetCriticalChance(action, source)
	roll := MakeActionRoll()
	success := float64(roll) < crit_chance
	ratio := 1.0
	if success {
		ratio = action.CritMod
	}

	fmt.Printf("crit_chance = %f, roll = %d\n", crit_chance, roll)

	return ChanceResult{
		Chance:  crit_chance,
		Roll:    roll,
		Success: success,
		Ratio:   ratio,
	}
}

func MakeAccuracyCheck(g Game, action ActionConfig, source ResolvedActor, target ResolvedActor, ignoreModifiers bool) ChanceResult {
	base_accuracy := GetAccuracy(g, source, target, ignoreModifiers)
	if action.Accuracy == nil {
		return ChanceResult{
			Success: true,
		}
	}

	accuracy := Round(base_accuracy * (float64(*action.Accuracy + source.ActionAccuracyOffset)))
	roll := MakeActionRoll()

	fmt.Printf("acc = %d, roll = %d\n", accuracy, roll)

	return ChanceResult{
		Chance:  float64(accuracy),
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
