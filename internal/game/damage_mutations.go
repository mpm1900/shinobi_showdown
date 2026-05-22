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
				result := MakeAccuracyCheck(g, action, source, target, false)
				if !result.Success {
					g.PushLog(NewLog(fmt.Sprintf("%s missed!", action.Name)))
					g.PushLog(NewLog(fmt.Sprintf("roll = %d, acc = %f", result.Roll, result.Chance)))
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
