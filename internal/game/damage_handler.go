package game

import (
	"fmt"

	"github.com/google/uuid"
)

type DamageConfig struct {
	Critical        float64
	Random          float64
	IgnoreProtect   bool
	IgnoreModifiers bool
	Repeat          bool
	RepeatMax       int
	OnSuccess       func(Game, Context, Context) []GameTransaction
	OnFailure       func(Game, Context, Context) []GameTransaction
}

type damageHandler struct {
	action          ActionConfig
	config          DamageConfig
	context         Context
	source          ResolvedActor
	resolvedTargets []ResolvedActor
	defense         ActorStat

	total  int
	totals []int

	repeats            int
	repeatTransactions []GameTransaction
	sideEffectTxs      []GameTransaction
}

func NewDamageConfig(critical float64, random float64) DamageConfig {
	return DamageConfig{
		Critical:        critical,
		Random:          random,
		IgnoreModifiers: critical > 1,
		IgnoreProtect:   false,
		Repeat:          false,
		RepeatMax:       0,
	}
}

func newDamageHandler(g Game, action ActionConfig, config DamageConfig, context Context, source ResolvedActor) *damageHandler {
	resolvedTargets := resolveTargets(g, context)
	return &damageHandler{
		action:             action,
		config:             config,
		context:            context,
		source:             source,
		resolvedTargets:    resolvedTargets,
		defense:            getDefenseStat(*action.Stat),
		totals:             make([]int, len(resolvedTargets)),
		repeatTransactions: make([]GameTransaction, 0),
		sideEffectTxs:      make([]GameTransaction, 0),
	}
}

func (e *damageHandler) run(g *Game) {
	for {
		missed := false

		for ti, target := range e.resolvedTargets {
			if target.HasJutsuImmunity(e.action.Jutsu) {
				log := MakeGameLog(fmt.Sprintf("$source$ was immune to %s", e.action.Jutsu), MakeContextForActor(target.Actor), 1)
				g.PushLog(log)
				continue
			}

			if e.resolveTargetHit(g, ti, target) {
				missed = true
			}
		}

		if !e.config.Repeat || missed {
			break
		}

		if e.config.RepeatMax <= 0 || e.config.RepeatMax > e.repeats+1 {
			e.repeats++
		} else {
			break
		}
	}

	e.buildSideEffects()
	e.commitTransactions(g)
}

func (e *damageHandler) resolveTargetHit(g *Game, targetIndex int, target ResolvedActor) bool {
	targetContext := e.context
	targetContext.TargetActorIDs = []uuid.UUID{target.ID}
	targetContext.TargetPositionIDs = []uuid.UUID{}

	if target.Protected && !e.config.IgnoreProtect {
		g.PushLog(MakeGameLog("$source$ was protected.", targetContext.WithSource(target.ID), 1))
		g.On(OnProtected, &targetContext)

		if e.config.OnFailure != nil {
			e.sideEffectTxs = append(e.sideEffectTxs, e.config.OnFailure(*g, e.context, targetContext)...)
		}

		return false
	}

	result := MakeAccuracyCheck(g, e.action, e.source, target, e.config.IgnoreModifiers)
	if !result.Success {
		if !e.config.Repeat || e.repeats == 0 {
			g.PushLog(NewLog(fmt.Sprintf("%s missed!", e.action.Name)))
			g.PushLog(NewLog(fmt.Sprintf("roll = %d, acc = %d", result.Roll, result.Chance)))
			g.On(OnMiss, &targetContext)
		}

		if e.config.OnFailure != nil {
			e.sideEffectTxs = append(e.sideEffectTxs, e.config.OnFailure(*g, e.context, targetContext)...)
		}

		return true
	}

	if e.config.OnSuccess != nil {
		e.sideEffectTxs = append(e.sideEffectTxs, e.config.OnSuccess(*g, e.context, targetContext)...)
	}

	damages := GetDamage(
		e.source,
		[]ResolvedActor{target},
		e.config.IgnoreModifiers,
		len(e.resolvedTargets),
		*e.action.Stat,
		e.defense,
		*e.action.Power,
		e.config.Critical,
		e.action.Nature,
		e.config.Random,
	)

	for _, damage := range damages {
		if e.config.Repeat {
			e.queueRepeatHit(g, target, damage)
		} else {
			e.applySingleHit(g, target, damage)
		}

		applied := clampDamage(damage)
		e.total += applied
		e.totals[targetIndex] += applied
	}

	e.handleNatureEffectiveLogging(g, target)

	return false
}

func (e *damageHandler) applySingleHit(g *Game, target ResolvedActor, damage int) {
	ApplyDamage(g, &e.source.ID, target, damage)

	if damage > 0 && e.context.SourceActorID != nil {
		ctx := MakeContextForActor(target.Actor).WithSource(*e.context.SourceActorID)
		if e.context.SourcePlayerID != nil {
			ctx = ctx.WithPlayer(*e.context.SourcePlayerID)
		}

		g.On(OnDamageReceive, &ctx)
		g.UpdateActor(target.ID, func(a Actor) Actor {
			a.HitCount++
			return a
		})

		if e.action.Stat != nil {
			stat := *e.action.Stat
			if stat == StatAttack || stat == StatDefense {
				g.On(OnDamagePhysical, &ctx)
			}
		}
	}

	if e.config.Critical > 1.0 {
		g.PushLog(MakeGameLog(fmt.Sprintf("Critical Hit! (x%f)", e.config.Critical), NewContext(), 1))
		g.On(OnCritical, &e.context)
	}
}

func (e *damageHandler) queueRepeatHit(g *Game, target ResolvedActor, damage int) {
	targetContext := e.context
	targetContext.TargetActorIDs = []uuid.UUID{target.ID}
	targetContext.TargetPositionIDs = []uuid.UUID{}

	repeatTx := MakeTransaction(PureDamage(damage, true), targetContext)

	log := NewLogContext(fmt.Sprintf("$action$ hit %d times.", e.repeats+1), e.context)
	logMux := AddLogs(log)
	logMux.Filter = TargetsAreOneAlive
	logTx := MakeTransaction(logMux, e.context)

	g.UpdateActor(target.ID, func(a Actor) Actor {
		a.HitCount++
		return a
	})

	if e.config.Critical > 1.0 {
		critlog := MakeGameLog(fmt.Sprintf("Critical Hit! (x%f)", e.config.Critical), NewContext(), 1)
		critlogMux := AddLogs(critlog)
		critlogMux.Filter = TargetsAreOneAlive
		critlogTx := MakeTransaction(critlogMux, e.context)
		triggerTx := RunTriggerTx(OnCritical, e.context)
		e.repeatTransactions = append(e.repeatTransactions, logTx, critlogTx, triggerTx)
	} else {
		e.repeatTransactions = append(e.repeatTransactions, logTx)
	}

	e.repeatTransactions = append(e.repeatTransactions, repeatTx)
}

func (e *damageHandler) handleNatureEffectiveLogging(g *Game, target ResolvedActor) {
	var natures []Nature
	if e.action.Nature != nil {
		natures = NATURES[*e.action.Nature]
	}

	natureResult := ResolveNatures(natures, e.source.NatureDamage, target.NatureResistance, target.Natures)
	if natureResult.Average >= NATURE_WEAKNESS_FULL {
		tctx := MakeContextForActor(target.Actor)
		tctx.SourceActorID = e.context.SourceActorID
		g.PushLog(MakeGameLog("Super effective!", tctx, 1))
		g.On(OnWeakness, &tctx)
	}
	if natureResult.Result <= NATURE_RESISTANCE_FULL {
		tctx := MakeContextForActor(target.Actor)
		tctx.SourceActorID = e.context.SourceActorID
		g.PushLog(MakeGameLog("Not very effective!", tctx, 1))
		g.On(OnResistance, &tctx)
	}
}

func (e *damageHandler) buildSideEffects() {
	if e.total == 0 || e.context.SourceActorID == nil {
		return
	}

	context := MakeContextForActor(e.source.Actor)
	if e.action.LifeSteal != nil && *e.action.LifeSteal > 0.0 {
		amount := Round(*e.action.LifeSteal * float64(e.total))
		healTx := MakeTransaction(PureHeal(amount), context)
		e.sideEffectTxs = append(e.sideEffectTxs, healTx)
	}

	if e.action.Recoil != nil {
		recoil := *e.action.Recoil * e.source.RecoilMultiplier
		if recoil > 0.0 {
			amount := Round(recoil * float64(e.total))
			recoilTx := MakeTransaction(PureDamage(amount, false), context)
			e.sideEffectTxs = append(e.sideEffectTxs, recoilTx)
		}
		if recoil < 0.0 {
			amount := Round(recoil * float64(e.total))
			recoilTx := MakeTransaction(PureHeal(amount*-1), context)
			e.sideEffectTxs = append(e.sideEffectTxs, recoilTx)
		}
	}

	for i, target := range e.resolvedTargets {
		if target.Reflect > 0.0 && *e.context.SourceActorID != target.ID {
			reflectDamage := int(target.Reflect * float64(e.totals[i]))
			reflectTx := MakeTransaction(PureDamage(reflectDamage, false), context)
			e.sideEffectTxs = append(e.sideEffectTxs, reflectTx)
		}
	}
}

func (e *damageHandler) commitTransactions(g *Game) {
	ordered := make([]GameTransaction, 0, len(e.repeatTransactions)+len(e.sideEffectTxs))
	ordered = append(ordered, e.repeatTransactions...)
	ordered = append(ordered, e.sideEffectTxs...)

	for i := len(ordered) - 1; i >= 0; i-- {
		g.JumpTransaction(ordered[i])
	}
}
