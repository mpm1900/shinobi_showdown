package game

import (
	"fmt"

	"github.com/google/uuid"
)

type DamageHit struct {
	Damage       int
	Critical     bool
	CiritcalMult float64
}

type DamageResult struct {
	Failed          bool
	Missed          bool
	Immune          bool
	IgnoreModifiers bool
	IgnoreProtect   bool
	Protected       bool
	Hits            []DamageHit
}

type DamageConfig struct {
	Random          float64
	IgnoreProtect   bool
	IgnoreModifiers bool
	Repeat          bool
	RepeatMax       int
	OnSuccess       func(Game, Context, Context) []GameTransaction
	OnFailure       func(Game, Context, Context) []GameTransaction
}

type DamageCore struct {
	ActionConfig ActionConfig
	DamageConfig DamageConfig

	Source  ResolvedActor
	Targets []ResolvedActor
	Context Context
	Results map[uuid.UUID]DamageResult
}

func NewDamageConfig(random float64) DamageConfig {
	return DamageConfig{
		Random:          random,
		IgnoreModifiers: false,
		IgnoreProtect:   false,
		Repeat:          false,
		RepeatMax:       0,
	}
}

func NewDamageCore(actionConfig ActionConfig, damageConfig DamageConfig, game Game, context Context) *DamageCore {
	s, ok := game.GetSource(context)
	if !ok {
		return nil
	}
	source := s.Resolve(game)
	targets := game.GetTargets(context)

	core := &DamageCore{
		ActionConfig: actionConfig,
		DamageConfig: damageConfig,
		Source:       source,
		Targets:      make([]ResolvedActor, len(targets)),
		Context:      context,
		Results:      make(map[uuid.UUID]DamageResult),
	}

	for i, t := range targets {
		core.Targets[i] = t.Resolve(game)
	}

	return core
}

func (dc *DamageCore) buildTargetHit(game Game, target ResolvedActor, result *DamageResult) {
	hit := DamageHit{
		Critical:     false,
		CiritcalMult: 1.0,
	}
	critical_result := MakeCriticalCheck(dc.ActionConfig, dc.Source)
	if critical_result.Success {
		hit.Critical = true
		hit.CiritcalMult = dc.ActionConfig.CritMod
	}

	hit.Damage = GetTargetDamage(
		dc.Source,
		target,
		result.IgnoreModifiers || hit.Critical,
		len(dc.Targets),
		*dc.ActionConfig.Stat,
		getDefenseStat(*dc.ActionConfig.Stat),
		*dc.ActionConfig.Power,
		hit.CiritcalMult,
		dc.ActionConfig.Nature,
		dc.DamageConfig.Random,
	)

	result.Hits = append(result.Hits, hit)
	if dc.DamageConfig.Repeat && len(result.Hits) < dc.DamageConfig.RepeatMax {
		accuracy_result := MakeAccuracyCheck(game, dc.ActionConfig, dc.Source, target, result.IgnoreModifiers)
		if accuracy_result.Success {
			dc.buildTargetHit(game, target, result)
		}
	}
}

func (dc *DamageCore) buildTargetResult(game Game, target ResolvedActor) {
	result := DamageResult{
		IgnoreModifiers: dc.DamageConfig.IgnoreModifiers,
		IgnoreProtect:   dc.DamageConfig.IgnoreProtect,
		Hits:            []DamageHit{},
	}

	defer func() {
		dc.Results[target.ID] = result
	}()

	if target.HasJutsuImmunity(dc.ActionConfig.Jutsu) {
		result.Failed = true
		result.Immune = true
	}

	if target.Protected && !result.IgnoreProtect {
		result.Failed = true
		result.Protected = true
		return
	}

	accuracy_result := MakeAccuracyCheck(game, dc.ActionConfig, dc.Source, target, result.IgnoreModifiers)
	if !accuracy_result.Success {
		result.Failed = true
		result.Missed = true
		return
	}

	result.Failed = false
	result.Immune = false
	result.Protected = false
	result.Missed = false

	dc.buildTargetHit(game, target, &result)
}

func (dc *DamageCore) BuildResults(game Game) {
	for _, target := range dc.Targets {
		dc.buildTargetResult(game, target)
	}
}

func log(str string, context Context) GameTransaction {
	gl := NewLogContext(str, context)
	mut := AddLogs(gl)
	return MakeTransaction(mut, context)
}
func (dc *DamageCore) logFailedResult(result DamageResult, target ResolvedActor) []GameTransaction {
	transactions := NewTransactionBuilder()

	target_context := MakeContextForActor(target.Actor)
	if result.Failed {
		if result.Immune {
			transactions.PushOne(log(fmt.Sprintf("$source$ was immune to %s", dc.ActionConfig.Jutsu), target_context))
		}
		if result.Protected {
			transactions.PushOne(log("$source$ was protected.", target_context))
		}

		if result.Missed {
			transactions.PushOne(log(fmt.Sprintf("%s missed!", dc.ActionConfig.Name), target_context))
		}
	}

	return transactions.Build()
}
func (dc *DamageCore) logHit(hit DamageHit, target ResolvedActor) []GameTransaction {
	transactions := NewTransactionBuilder()
	target_context := MakeContextForActor(target.Actor)

	if hit.Critical {
		transactions.PushOne(log(fmt.Sprintf("Critical Hit! (x%f)", hit.CiritcalMult), target_context))
	}

	if dc.ActionConfig.Nature != nil {
		natureResult := ResolveNatures(dc.ActionConfig.Nature.GetNatures(), dc.Source.NatureDamage, target.NatureResistance, target.Natures)
		if natureResult.Average >= NATURE_WEAKNESS_FULL {
			transactions.PushOne(log("Super effective!", target_context))
		}
		if natureResult.Result <= NATURE_RESISTANCE_FULL {
			transactions.PushOne(log("Not very effective!", target_context))
		}
	}

	return transactions.Build()
}

func (dc *DamageCore) resolveFailedResult(game Game, result DamageResult, target ResolvedActor) []GameTransaction {
	transactions := NewTransactionBuilder()

	transactions.Push(dc.logFailedResult(result, target))
	triggers := GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			if result.Protected {
				g.On(OnProtected, &context)
			}
			if result.Missed {
				g.On(OnMiss, &context)
			}
			return g
		},
	}
	target_context := MakeContextForActor(target.Actor).WithSource(dc.Source.ID).WithPlayer(dc.Source.PlayerID)
	transactions.PushOne(MakeTransaction(triggers, target_context))

	if dc.DamageConfig.OnFailure != nil {
		transactions.Push(dc.DamageConfig.OnFailure(game, dc.Context, target_context))
	}

	return transactions.Build()
}

func (dc *DamageCore) resolveSuccessResult(game Game, result DamageResult, target ResolvedActor) []GameTransaction {
	transactions := NewTransactionBuilder()

	target_context := MakeContextForActor(target.Actor).WithSource(dc.Source.ID).WithPlayer(dc.Source.PlayerID)
	for index, hit := range result.Hits {
		if hit.Damage == 0 {
			continue
		}

		damage_mut := GameMutation{
			Delta: func(p Game, g Game, context Context) Game {
				ApplyDamage(&g, &dc.Source.Actor, target, hit.Damage)
				if hit.Damage > 0 {
					g.On(OnDamageReceive, &context)
					g.UpdateActor(target.ID, func(a Actor) Actor {
						a.HitCount++
						return a
					})

					if dc.ActionConfig.Stat != nil {
						stat := *dc.ActionConfig.Stat
						if stat == StatAttack || stat == StatDefense {
							g.On(OnDamagePhysical, &context)
						}
					}
				}
				return g
			},
		}

		transactions.PushOne(MakeTransaction(damage_mut, target_context))

		if len(result.Hits) > 1 {
			repeat_log := log(fmt.Sprintf("$action$ hit %d time.", index+1), dc.Context)
			if index > 0 {
				repeat_log = log(fmt.Sprintf("$action$ hit %d times.", index+1), dc.Context)
			}
			repeat_log.Mutation.Filter = TargetsAreOneAlive
			transactions.PushOne(repeat_log)
		}

		transactions.Push(dc.logHit(hit, target))

		source_context := MakeContextForActor(dc.Source.Actor)
		damage_amount := min(float64(hit.Damage), float64(target.Stats[StatHP]))
		if dc.ActionConfig.LifeSteal != nil && *dc.ActionConfig.LifeSteal > 0.0 {
			amount := Round(*dc.ActionConfig.LifeSteal * damage_amount)
			healTx := MakeTransaction(PureHeal(amount), source_context)
			transactions.PushOne(healTx)
		}

		if dc.ActionConfig.Recoil != nil {
			recoil := *dc.ActionConfig.Recoil * dc.Source.RecoilMultiplier
			amount := Round(recoil * damage_amount)
			if recoil > 0.0 {
				recoilTx := MakeTransaction(PureDamage(amount, false), source_context)
				transactions.PushOne(recoilTx)
			}
			if recoil < 0.0 {
				recoilTx := MakeTransaction(PureHeal(amount*-1), source_context)
				transactions.PushOne(recoilTx)
			}
		}

		if target.Reflect > 0.0 {
			reflectDamage := Round(target.Reflect * damage_amount)
			reflectTx := MakeTransaction(PureDamage(reflectDamage, false), source_context)
			transactions.PushOne(reflectTx)
		}
	}

	if dc.DamageConfig.OnSuccess != nil {
		transactions.Push(dc.DamageConfig.OnSuccess(game, dc.Context, target_context))
	}

	return transactions.Build()
}

func (dc *DamageCore) ResolveResults(game Game) []GameTransaction {
	transactions := NewTransactionBuilder()

	for _, target := range dc.Targets {
		if result, ok := dc.Results[target.ID]; ok {
			if result.Failed {
				transactions.Push(dc.resolveFailedResult(game, result, target))
			}

			if !result.Failed {
				transactions.Push(dc.resolveSuccessResult(game, result, target))
			}
		}
	}

	return transactions.Build()
}

func (dc *DamageCore) Run(game Game) []GameTransaction {
	dc.BuildResults(game)
	return dc.ResolveResults(game)
}

func DamageCoreMutation(actionConfig ActionConfig, damageConfig DamageConfig) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			core := NewDamageCore(actionConfig, damageConfig, g, context)
			transactions := core.Run(g)
			g.JumpTransactions(transactions)

			return g
		},
	}
}

func ResolveDamageCore(actionConfig ActionConfig, damageConfig DamageConfig, game Game, context Context) []GameTransaction {
	core := NewDamageCore(actionConfig, damageConfig, game, context)
	return core.Run(game)
}
