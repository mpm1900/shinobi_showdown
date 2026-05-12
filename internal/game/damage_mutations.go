package game

import (
	"fmt"

	"github.com/google/uuid"
)

func clampDamage(damage int) int {
	if damage < 0 {
		return 0
	}
	return damage
}

func resolveTargets(g Game, context Context) []ResolvedActor {
	targets := g.GetTargets(context)
	resolved := make([]ResolvedActor, len(targets))
	for i, t := range targets {
		resolved[i] = t.Resolve(g)
	}
	return resolved
}

func getDefenseStat(a ActorStat) ActorStat {
	if a == StatChakraAttack {
		return StatChakraDefense
	}
	if a == StatChakraDefense {
		return StatChakraDefense
	}
	return StatDefense
}

func SetDamage(damage int) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			targets := g.GetTargets(context)
			for _, t := range targets {
				g.UpdateActor(t.ID, func(a Actor) Actor {
					a.Damage = damage
					return a
				})
			}
			return g
		},
	}
}

// returns if target is still alive after
func ApplyDamageWith(g *Game, source_ID *uuid.UUID, target ResolvedActor, damage int, updater func(Actor) Actor) bool {
	alive := target.Alive
	if !target.Alive {
		return alive
	}

	logCtx := NewContext()
	logCtx.ParentActorID = &target.ID
	logCtx.SourceActorID = &target.ID

	hp := target.Stats[StatHP]

	g.UpdateActor(target.ID, func(a Actor) Actor {
		if a.Summon != nil && a.Summon.Proxy && a.Summon.Alive {
			summonHP := a.Summon.Stats[StatHP]
			a.Summon.Damage += clampDamage(damage)
			a.Summon.Alive = summonHP > a.Summon.Damage
			g.PushLog(MakeGameLog("$source$'s summon took the attack.", logCtx, 1))
		} else {
			a.Damage += damage
			if source_ID != nil {
				a.LastReceivedDamage[*source_ID] = clampDamage(damage)
			}
			ratio := min(int(float64(damage)*100/float64(hp)), 100)
			if ratio > 0 {
				g.PushLog(MakeGameLog(fmt.Sprintf("$source$ lost %d%% HP.", ratio), logCtx, 1))
			} else {
				g.PushLog(MakeGameLog(fmt.Sprintf("$source$ gained %d%% HP.", ratio*-1), logCtx, 1))
			}

			if target.Immortal && hp <= a.Damage {
				g.PushLog(MakeGameLog("$source$'s survived a fatal attack.", logCtx, 1))
				a.Damage = hp - 1
				g.On(OnImmortalSave, &logCtx)
			}

			a.Alive = hp > a.Damage
			alive = a.Alive
		}

		if updater == nil {
			return a
		}
		return updater(a)
	})

	if !alive {
		deathContext := NewContext().WithSource(target.ID)
		g.On(OnDeath, &deathContext)
		if source_ID != nil {
			killContext := NewContext().WithSource(*source_ID).WithTargetIDs([]uuid.UUID{target.ID})
			g.On(OnKill, &killContext)
		}
	}

	return alive
}

// returns if target is still alive after
func ApplyDamage(g *Game, source_ID *uuid.UUID, target ResolvedActor, damage int) bool {
	return ApplyDamageWith(g, source_ID, target, damage, nil)
}

func PureDamageWith(damage int, trigger bool, updater func(Actor) Actor) GameMutation {
	return GameMutation{
		Filter: TargetsAreOneAlive,
		Delta: func(p Game, g Game, context Context) Game {
			targets := g.GetTargets(context)
			for _, t := range targets {
				target := t.Resolve(g)
				ApplyDamageWith(&g, context.SourceActorID, target, damage, updater)
				if trigger && damage > 0 {
					g.On(OnDamageReceive, &context)
				}
			}
			return g
		},
	}
}

func PureDamage(damage int, trigger bool) GameMutation {
	return PureDamageWith(damage, trigger, nil)
}

func RatioDamageWith(ratio float64, updater func(Actor) Actor) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			targets := g.GetTargets(context)
			for _, t := range targets {
				target := t.Resolve(g)
				damage := Round(float64(target.Stats[StatHP]) * ratio)
				ApplyDamageWith(&g, context.SourceActorID, target, damage, updater)
			}
			return g
		},
	}
}

func RatioDamage(ratio float64) GameMutation {
	return RatioDamageWith(ratio, nil)
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

	e.handleNatureEffectiveness(g, target)

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
func (e *damageHandler) handleNatureEffectiveness(g *Game, target ResolvedActor) {
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

	if e.action.Recoil != nil && *e.action.Recoil > 0.0 {
		amount := Round(*e.action.Recoil * float64(e.total))
		recoilTx := MakeTransaction(PureDamage(amount, false), context)
		e.sideEffectTxs = append(e.sideEffectTxs, recoilTx)
	}
	if e.action.Recoil != nil && *e.action.Recoil < 0.0 {
		amount := Round(*e.action.Recoil * float64(e.total))
		recoilTx := MakeTransaction(PureHeal(amount*-1), context)
		e.sideEffectTxs = append(e.sideEffectTxs, recoilTx)
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

func NewDamage(action ActionConfig, config DamageConfig) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			s, ok := g.GetSource(context)
			if !ok || action.Stat == nil || action.Power == nil {
				return g
			}

			source := s.Resolve(g)
			exec := newDamageHandler(g, action, config, context, source)
			exec.run(&g)

			return g
		},
	}
}

func MakeDamageTransactions(context Context, damages ...GameMutation) []GameTransaction {
	var transactions []GameTransaction
	for _, damage := range damages {
		transactions = append(
			transactions,
			MakeTransaction(
				damage,
				context,
			),
		)
	}
	return transactions
}

func ApplyHealRawWith(g *Game, targetID uuid.UUID, amount int, updater func(Actor) Actor) int {
	g.UpdateActor(targetID, func(a Actor) Actor {
		if !a.Alive {
			amount = 0
			return a
		}

		healed := min(amount, a.Damage)
		a.Damage -= healed
		amount = healed

		if updater == nil {
			return a
		}
		return updater(a)
	})

	t, ok := g.GetActorByID(targetID)
	if !ok {
		return amount
	}

	target := t.Resolve(*g)
	hp := target.Stats[StatHP]
	logCtx := MakeContextForActor(target.Actor)
	ratio := int(float64(amount) * 100 / float64(hp))
	g.PushLog(MakeGameLog(fmt.Sprintf("$source$ gained %d%% HP.", ratio), logCtx, 1))

	return amount
}
func ApplyHealRaw(g *Game, targetID uuid.UUID, amount int) int {
	return ApplyHealRawWith(g, targetID, amount, nil)
}
func ApplyHealRatioWith(g *Game, target ResolvedActor, ratio float64, updater func(Actor) Actor) int {
	amount := Round(float64(target.Stats[StatHP]) * ratio)
	return ApplyHealRawWith(g, target.ID, amount, updater)
}
func ApplyHealRatio(g *Game, target ResolvedActor, ratio float64) int {
	return ApplyHealRatioWith(g, target, ratio, nil)
}
func RatioHeal(ratio float64) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			for _, target := range resolveTargets(g, context) {
				ApplyHealRatio(&g, target, ratio)
			}
			return g
		},
	}
}
func PureHeal(amount int) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			targets := g.GetTargets(context)
			for _, target := range targets {
				ApplyHealRaw(&g, target.ID, amount)
			}
			return g
		},
	}
}
func NewHeal(action ActionConfig, ratio float64) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			s, ok := g.GetSource(context)
			if !ok {
				return g
			}

			source := s.Resolve(g)
			for _, target := range resolveTargets(g, context) {
				result := MakeAccuracyCheck(&g, action, source, target, false)
				if !result.Success {
					g.PushLog(NewLog(fmt.Sprintf("%s missed!", action.Name)))
					g.PushLog(NewLog(fmt.Sprintf("roll = %d, acc = %d", result.Roll, result.Chance)))
					continue
				}
				ApplyHealRatio(&g, target, ratio)
			}

			return g
		},
	}
}

func SetHpSourceToTargets() GameMutation {
	return GameMutation{
		Delta: func(p, g Game, context Context) Game {
			s, ok := g.GetSource(context)
			if !ok {
				return g
			}

			source := s.Resolve(g)
			source_hp := source.Stats[StatHP] - source.Damage
			for _, target := range resolveTargets(g, context) {
				target_hp := target.Stats[StatHP] - target.Damage
				diff := source_hp - target_hp
				if diff > 0 {
					ApplyHealRaw(&g, target.ID, diff)
				} else {
					ApplyDamage(&g, &source.ID, target, diff*-1)
				}
			}

			return g
		},
	}
}
