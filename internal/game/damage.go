package game

import (
	"fmt"
	"maps"
	"math/rand/v2"
	"slices"
)

type DamageTerms struct {
	Attack   int
	Critical float64
	Defense  int
	Level    int
	Nature   NatureResult
	Offset   int
	Other    float64
	Power    int
	Random   float64
	STAB     float64
	Targets  float64
}

func RandomDamageFactor() float64 {
	roll := rand.IntN(15) + 86
	value := float64(roll) / 100.0
	return value
}

func DamageEquation(terms DamageTerms) int {
	pow_ad := float64(terms.Power) * float64(terms.Attack) / float64(terms.Defense)
	level_mod := float64(2*terms.Level)/5 + 2
	base := (pow_ad*level_mod)/50 + 2
	mods := terms.Critical * terms.Nature.Result * terms.STAB * terms.Targets * terms.Random * terms.Other
	raw := (base * mods)
	fmt.Println("DAMAGE REPORT:")
	fmt.Printf("A/D POW = %f\n \t%d * %d / %d\n", pow_ad, terms.Power, terms.Attack, terms.Defense)
	fmt.Printf("LEVEL = %f\n \t(2 * %d / 5) + 2\n", level_mod, terms.Level)
	fmt.Printf("BASE = %f\n \t(%f * %f) / 50 + 2\n", base, pow_ad, level_mod)
	fmt.Printf("MODS = %f\n \tcrit=%f\n \tnature=%f\n \tnature_base=%f\n \tstab=%f\n \ttargets=%f\n \trandom=%f\n \tother=%f\n", mods, terms.Critical, terms.Nature.Result, terms.Nature.Average, terms.STAB, terms.Targets, terms.Random, terms.Other)
	fmt.Printf(
		"TOTAL = %f\n \t([%d * %d / %d] * [(2 * %d / 5) + 2] / 50 + 2) * (%f * %f * %f * %f * %f * %f) = %f \n",
		raw, terms.Power, terms.Attack, terms.Defense, terms.Level, terms.Critical, terms.Nature.Result, terms.STAB, terms.Targets, terms.Random, terms.Other, raw,
	)
	return Round(raw) + terms.Offset
}

func GetStabModifier(source ResolvedActor, nature *NatureSet) float64 {
	if nature == nil {
		return 1.00
	}

	natures := slices.Collect(maps.Keys(source.Natures))
	index := slices.IndexFunc(natures, func(n NatureSet) bool {
		return n == *nature
	})

	if index == -1 {
		return 1.00
	}

	return source.StabMultiplier
}

func HasDebuff(r ResolvedActor, stat AttackStat) bool {
	damage_mult, ok := r.DamageMultipliers[stat]
	if !ok {
		return false
	}
	return damage_mult < 1 ||
		r.Stages[ActorStat(stat)] < 0 ||
		r.Stages[StatAccuracy] < 0
}
func HasBuff(r ResolvedActor, attack AttackStat, defense DefenseStat) bool {
	damage_reduction, ok := r.DamageReduction[attack]
	if !ok {
		return false
	}
	return damage_reduction > 1 ||
		r.Stages[ActorStat(defense)] > 0 ||
		r.Stages[StatEvasion] > 0
}

func GetTargetDamage(
	source ResolvedActor,
	target ResolvedActor,
	ignoreModifiers bool,
	totalTargets int,
	attack ActorStat,
	defense ActorStat,
	power int,
	critical float64,
	nature *NatureSet,
	random float64,
) int {
	if power == 0 {
		return 0
	}

	a_base := float64(source.Stats[attack])
	// only ignore mods if attack is weakened
	if ignoreModifiers && HasDebuff(source, AttackStat(attack)) {
		a_base = float64(source.UnmodifiedStats[attack])
	}
	a_mod := 1.0
	attack_value := Round(a_base * a_mod)

	targets_mod := 1.0
	if totalTargets > 1 {
		targets_mod = 0.75
	}

	d_base := float64(target.Stats[defense])
	// only ignore mods if defense is strengthened
	if ignoreModifiers && HasBuff(target, AttackStat(attack), DefenseStat(defense)) {
		d_base = float64(target.UnmodifiedStats[defense])
	}
	d_mod := 1.0
	defense_value := Round(d_base * d_mod)

	var natures []Nature
	if nature != nil {
		natures = NATURES[*nature]
	}
	nature_mod := ResolveNatures(natures, source.NatureDamage, target.NatureResistance, target.Natures)
	stab_mod := GetStabModifier(source, nature)
	damage_mult, ok := source.DamageMultipliers[AttackStat(attack)]
	if !ok {
		damage_mult = 1
	}
	damage_reduction, ok := target.DamageReduction[AttackStat(attack)]
	if !ok {
		damage_reduction = 1
	}

	return DamageEquation(DamageTerms{
		Attack:   attack_value,
		Critical: critical,
		Defense:  defense_value,
		Level:    source.Level,
		Nature:   nature_mod,
		Offset:   0,
		Other:    damage_mult * damage_reduction,
		Power:    Round(float64(power) * source.PowerMultiplier),
		Random:   random,
		STAB:     stab_mod,
		Targets:  targets_mod,
	})
}
