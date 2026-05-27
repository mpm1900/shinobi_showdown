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

type DamageCoreConfig struct {
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
	DamageConfig DamageCoreConfig

	Source  ResolvedActor
	Targets []ResolvedActor
	Context Context
	Results map[uuid.UUID]DamageResult
}

func NewDamageConfig(random float64) DamageCoreConfig {
	return DamageCoreConfig{
		Random:          random,
		IgnoreModifiers: false,
		IgnoreProtect:   false,
		Repeat:          false,
		RepeatMax:       0,
	}
}

func NewDamageCore(actionConfig ActionConfig, damageConfig DamageCoreConfig, game Game, context Context) *DamageCore {
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
	var transactions []GameTransaction

	target_context := MakeContextForActor(target.Actor)
	if result.Failed {
		if result.Immune {
			transactions = append(transactions, log(fmt.Sprintf("$source$ was immune to %s", dc.ActionConfig.Jutsu), target_context))
		}
		if result.Protected {
			transactions = append(transactions, log("$source$ was protected.", target_context))
		}

		if result.Missed {
			transactions = append(transactions, log(fmt.Sprintf("%s missed!", dc.ActionConfig.Name), target_context))
		}
	}

	return transactions
}
func (dc *DamageCore) logHit(hit DamageHit, target ResolvedActor) []GameTransaction {
	var transactions []GameTransaction
	target_context := MakeContextForActor(target.Actor)

	if hit.Critical {
		transactions = append(transactions, log(fmt.Sprintf("Critical Hit! (x%f)", hit.CiritcalMult), target_context))
	}

	var natures []Nature
	if dc.ActionConfig.Nature != nil {
		natures = NATURES[*dc.ActionConfig.Nature]
	}
	natureResult := ResolveNatures(natures, dc.Source.NatureDamage, target.NatureResistance, target.Natures)
	if natureResult.Average >= NATURE_WEAKNESS_FULL {
		transactions = append(transactions, log("Super effective!", target_context))
	}
	if natureResult.Result <= NATURE_RESISTANCE_FULL {
		transactions = append(transactions, log("Not very effective!", target_context))
	}

	return transactions
}

func (dc *DamageCore) ResolveResults(game *Game) []GameTransaction {
	var transactions []GameTransaction

	for _, target := range dc.Targets {
		target_context := MakeContextForActor(target.Actor).WithSource(dc.Source.ID).WithPlayer(dc.Source.PlayerID)
		if result, ok := dc.Results[target.ID]; ok {
			if result.Failed {
				transactions = append(transactions, dc.logFailedResult(result, target)...)
			}

			if !result.Failed {
				for index, hit := range result.Hits {
					transactions = append(transactions, dc.logHit(hit, target)...)

					damage_mut := GameMutation{
						Delta: func(p Game, g Game, context Context) Game {
							ApplyDamage(&g, &dc.Source.ID, target, hit.Damage)
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
					transactions = append(transactions, MakeTransaction(damage_mut, target_context))

					if len(result.Hits) > 1 {
						repeat_log := log(fmt.Sprintf("$action$ hit %d time.", index+1), dc.Context)
						if index > 0 {
							repeat_log = log(fmt.Sprintf("$action$ hit %d times.", index+1), dc.Context)
						}
						repeat_log.Mutation.Filter = TargetsAreOneAlive
						transactions = append(transactions, repeat_log)
					}
				}
			}
		}
	}

	return transactions
}

func (dc *DamageCore) ResolveSideEffects(game *Game) []GameTransaction {
	var transactions []GameTransaction

	for _, target := range dc.Targets {
		target_context := MakeContextForActor(target.Actor).WithSource(dc.Source.ID).WithPlayer(dc.Source.PlayerID)
		if result, ok := dc.Results[target.ID]; ok {
			if result.Failed {
				if result.Protected {
					game.On(OnProtected, &target_context)
				}
				if result.Missed {
					game.On(OnMiss, &target_context)
				}

				if dc.DamageConfig.OnFailure != nil {
					transactions = append(transactions, dc.DamageConfig.OnFailure(*game, dc.Context, target_context)...)
				}
			}

			if !result.Failed {
				for _, hit := range result.Hits {
					if hit.Damage == 0 {
						continue
					}

					source_context := MakeContextForActor(dc.Source.Actor)
					if dc.ActionConfig.LifeSteal != nil && *dc.ActionConfig.LifeSteal > 0.0 {
						amount := Round(*dc.ActionConfig.LifeSteal * float64(hit.Damage))
						healTx := MakeTransaction(PureHeal(amount), source_context)
						transactions = append(transactions, healTx)
					}

					if dc.ActionConfig.Recoil != nil {
						recoil := *dc.ActionConfig.Recoil * dc.Source.RecoilMultiplier
						amount := Round(recoil * float64(hit.Damage))
						if recoil > 0.0 {
							recoilTx := MakeTransaction(PureDamage(amount, false), source_context)
							transactions = append(transactions, recoilTx)
						}
						if recoil < 0.0 {
							recoilTx := MakeTransaction(PureHeal(amount*-1), source_context)
							transactions = append(transactions, recoilTx)
						}
					}

					if target.Reflect > 0.0 {
						reflectDamage := int(target.Reflect * float64(hit.Damage))
						reflectTx := MakeTransaction(PureDamage(reflectDamage, false), source_context)
						transactions = append(transactions, reflectTx)
					}
				}
				if dc.DamageConfig.OnSuccess != nil {
					transactions = append(transactions, dc.DamageConfig.OnSuccess(*game, dc.Context, target_context)...)
				}
			}
		}
	}

	return transactions
}

func (dc *DamageCore) Run(game *Game) []GameTransaction {
	var transactions []GameTransaction
	dc.BuildResults(*game)
	transactions = append(transactions, dc.ResolveResults(game)...)
	transactions = append(transactions, dc.ResolveSideEffects(game)...)

	return transactions
}

func DamageCoreMutation(actionConfig ActionConfig, damageConfig DamageCoreConfig) GameMutation {
	return GameMutation{
		Delta: func(p Game, g Game, context Context) Game {
			core := NewDamageCore(actionConfig, damageConfig, g, context)
			transactions := core.Run(&g)
			g.JumpTransactions(transactions)

			return g
		},
	}
}
