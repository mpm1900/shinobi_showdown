package game

import (
	"maps"
	"slices"
)

type Nature string

const (
	NatureFire      Nature = "fire"
	NatureWind      Nature = "wind"
	NatureLightning Nature = "lightning"
	NatureEarth     Nature = "earth"
	NatureWater     Nature = "water"
	NatureYin       Nature = "yin"
	NatureYang      Nature = "yang"
	NatureTai       Nature = "tai"
)

type NatureSet string

const (
	NsTai       NatureSet = "tai"
	NsFire      NatureSet = NatureSet(NatureFire)
	NsWind      NatureSet = NatureSet(NatureWind)
	NsLightning NatureSet = NatureSet(NatureLightning)
	NsEarth     NatureSet = NatureSet(NatureEarth)
	NsWater     NatureSet = NatureSet(NatureWater)
	NsYin       NatureSet = NatureSet(NatureYin)
	NsYang      NatureSet = NatureSet(NatureYang)

	NsScorch    NatureSet = "scorch"
	NsLava      NatureSet = "lava"
	NsBoil      NatureSet = "boil"
	NsGale      NatureSet = "gale"
	NsMagnet    NatureSet = "magnet"
	NsIce       NatureSet = "ice"
	NsExplosion NatureSet = "explosion"
	NsStorm     NatureSet = "storm"
	NsWood      NatureSet = "wood"
	NsYinYang   NatureSet = "yinyang"
	NsParticle  NatureSet = "particle"
	NsPure      NatureSet = "pure"
	NsJashin    NatureSet = "jashin"
)

var NATURES = map[NatureSet][]Nature{
	NsTai:       {NatureTai},
	NsFire:      {NatureFire},
	NsWind:      {NatureWind},
	NsLightning: {NatureLightning},
	NsEarth:     {NatureEarth},
	NsWater:     {NatureWater},
	NsYin:       {NatureYin},
	NsYang:      {NatureYang},
	NsPure:      {},

	NsScorch:    {NatureFire, NatureWind},
	NsLava:      {NatureFire, NatureEarth},
	NsBoil:      {NatureFire, NatureWater},
	NsGale:      {NatureWind, NatureLightning},
	NsMagnet:    {NatureWind, NatureEarth},
	NsIce:       {NatureWind, NatureWater},
	NsExplosion: {NatureLightning, NatureEarth},
	NsStorm:     {NatureLightning, NatureWater},
	NsWood:      {NatureEarth, NatureWater},
	NsYinYang:   {NatureYin, NatureYang},
	NsParticle:  {NatureFire, NatureEarth, NatureWind},
	NsJashin:    {},
}

type NatureSetValue = map[Nature]float64

var NATURE_EFFECTIVENESS = map[Nature]NatureSetValue{
	NatureFire: {
		NatureFire:  0.8,
		NatureWind:  2,
		NatureWater: 0.5,
	},
	NatureWind: {
		NatureFire:      0.5,
		NatureWind:      0.8,
		NatureLightning: 2,
	},
	NatureLightning: {
		NatureWind:      0.5,
		NatureLightning: 0.8,
		NatureEarth:     2,
	},
	NatureEarth: {
		NatureFire:      1.25,
		NatureLightning: 0.5,
		NatureEarth:     0.8,
		NatureWater:     2,
	},
	NatureWater: {
		NatureFire:  2,
		NatureEarth: 0.5,
		NatureWater: 0.8,
	},
	NatureYin: {
		NatureYin:  1.25,
		NatureYang: 0.8,
	},
	NatureYang: {
		NatureYang: 1.25,
		NatureYin:  0.8,
	},
	NatureTai: {
		NatureTai: 0.8,
	},
}

func GetEffectiveRatio(action Nature, target []Nature) float64 {
	set_value := NATURE_EFFECTIVENESS[action]
	total := 1.0

	if len(target) == 0 {
		return 1.0
	}
	for _, t_nature := range target {
		nature_value, ok := set_value[t_nature]
		if !ok {
			total *= 1.0
			continue
		}

		total *= nature_value
	}

	return total
}

func NewNatureSetValues() map[Nature]float64 {
	return map[Nature]float64{
		NatureFire:      1.00,
		NatureWind:      1.00,
		NatureLightning: 1.00,
		NatureEarth:     1.00,
		NatureWater:     1.00,
		NatureYin:       1.00,
		NatureYang:      1.00,
	}
}

func MapNatures(keys []NatureSet) map[NatureSet][]Nature {
	natures := make(map[NatureSet][]Nature)
	for _, key := range keys {
		natures[key] = NATURES[key]
	}
	return natures
}

func CollectNatures(natures map[NatureSet][]Nature) []Nature {
	result := make(map[Nature]struct{})
	for _, group := range natures {
		for _, nature := range group {
			result[nature] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(result))
}

type NatureResult struct {
	Base    float64
	Average float64
	Mult    float64
	Result  float64
}

func ResolveNatures(
	input []Nature,
	damages map[Nature]float64,
	resistances map[Nature]float64,
	natures_map map[NatureSet][]Nature,
) NatureResult {
	natures := CollectNatures(natures_map)
	base_effectiveness := 1.0

	if len(input) > 0 {
		for _, moveNature := range input {
			base_effectiveness *= GetEffectiveRatio(moveNature, natures)
		}
	}

	var avg_effectiveness float64
	if len(input) == 0 {
		avg_effectiveness = 1
	} else {
		avg_effectiveness = base_effectiveness
	}

	mult := 1.0
	for _, nature := range input {
		res := resistances[nature]

		if res == 0 {
			mult = 0
			break
		}

		mult = mult * damages[nature] / res
	}

	result := 0.0
	if mult != 0 && avg_effectiveness != 0 {
		result = max(0, mult+avg_effectiveness-1.0)
	}

	return NatureResult{
		Base:    base_effectiveness,
		Average: avg_effectiveness,
		Mult:    mult,
		Result:  result,
	}
}
