import { entries } from '../utils'

type Nature =
  | 'tai'
  | 'pure'
  | 'fire'
  | 'wind'
  | 'lightning'
  | 'earth'
  | 'water'
  | 'yang'
  | 'yin'
type NatureSet =
  | Nature
  | 'scorch'
  | 'lava'
  | 'boil'
  | 'magnet'
  | 'ice'
  | 'explosion'
  | 'storm'
  | 'wood'
  | 'yinyang'
  | 'particle'
  | 'jashin'

type NatureEffectiveness = Record<Nature, Partial<Record<Nature, number>>>
const natureEffectiveness: NatureEffectiveness = {
  fire: {
    wind: 2,
    water: 0.5,
  },
  wind: {
    fire: 0.5,
    lightning: 2,
  },
  lightning: {
    wind: 0.5,
    earth: 2,
  },
  earth: {
    lightning: 0.5,
    water: 2,
  },
  water: {
    fire: 2,
    earth: 0.5,
  },
  tai: {},
  yang: {
    yin: 1.25,
    yang: 0.8,
  },
  yin: {
    yang: 1.25,
    yin: 0.8,
  },
  pure: {},
}

const natureSetMap: Record<NatureSet, Array<Nature>> = {
  tai: ['tai'],
  pure: ['pure'],
  fire: ['fire'],
  wind: ['wind'],
  lightning: ['lightning'],
  earth: ['earth'],
  water: ['water'],
  yang: ['yang'],
  yin: ['yin'],
  scorch: ['fire', 'wind'],
  lava: ['fire', 'earth'],
  boil: ['fire', 'water'],
  magnet: ['wind', 'earth'],
  ice: ['wind', 'water'],
  explosion: ['earth', 'lightning'],
  storm: ['lightning', 'water'],
  wood: ['earth', 'water'],
  yinyang: ['yin', 'yang'],
  particle: ['fire', 'earth', 'lightning'],
  jashin: [],
}

const natureNames: Partial<Record<NatureSet, string>> = {
  tai: '体',
  pure: '纯',
  fire: '火',
  wind: '風',
  lightning: '雷',
  earth: '土',
  water: '水',
  yin: '陰',
  yang: '陽',
  magnet: '磁',
  ice: '氷',
  explosion: '爆',
  storm: '嵐',
  wood: '木',
  yinyang: '陰陽',
  particle: '塵',
  jashin: '邪',
}

const natureIndexes: Record<NatureSet, number> = {
  tai: -1,
  fire: 0,
  wind: 1,
  lightning: 2,
  earth: 3,
  water: 4,
  yin: 5,
  yang: 6,
  scorch: 8,
  lava: 9,
  boil: 10,
  magnet: 12,
  ice: 13,
  explosion: 14,
  storm: 15,
  wood: 16,
  yinyang: 17,
  particle: 18,
  pure: 19,
  jashin: 20,
}

function getWeakness(...natures: Array<NatureSet>): Array<Nature> {
  const list = natures.flatMap((nature) => {
    const base = natureSetMap[nature]
    return base
      .flatMap((n) => {
        const kv = entries(natureEffectiveness[n])
        return kv.filter(([_, v]) => v == 0.5).map(([k]) => k)
      })
      .filter((n) => n !== undefined)
  })
  return Array.from(new Set(list))
}

function getResistance(...natures: Array<NatureSet>): Array<Nature> {
  const list = natures.flatMap((nature) => {
    const base = natureSetMap[nature]
    return base
      .flatMap((n) => {
        const kv = entries(natureEffectiveness[n])
        return kv.filter(([_, v]) => v == 2).map(([k]) => k)
      })
      .filter((n) => n !== undefined)
  })
  return Array.from(new Set(list))
}

function getEffectiveness(action: NatureSet, target_ns: NatureSet[]): number {
  const action_natures = natureSetMap[action]
  if (action_natures.length === 0) {
    return 1
  }

  const target_natures = Array.from(
    new Set(target_ns.flatMap((target) => natureSetMap[target]))
  )
  if (target_natures.length === 0) {
    return 1
  }

  const base = action_natures.reduce((total, action_nature) => {
    const effectiveness = natureEffectiveness[action_nature]
    const compounded = target_natures.reduce((acc, target_nature) => {
      return acc * (effectiveness[target_nature] ?? 1)
    }, 1)
    return total * compounded
  }, 1)

  return base
}

export type { Nature, NatureSet }
export {
  natureNames,
  natureIndexes,
  natureSetMap,
  natureEffectiveness,
  getWeakness,
  getResistance,
  getEffectiveness,
}
