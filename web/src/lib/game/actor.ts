import { clamp } from '#/utils/clamp'
import type { Action } from './action'
import type { Modifier } from './modifier'
import type { Nature, NatureSet } from './nature'

type ActorAttackStat = 'attack' | 'chakra_attack'
type ActorDefenseStat = 'defense' | 'chakra_defense'

type ActorNatureStat =
  | ActorAttackStat
  | ActorDefenseStat
  | 'accuracy'
  | 'evasion'
  | 'speed'

type ActorBaseStat = ActorNatureStat | 'hp' | 'stamina'

type ActorStats<T> = Record<ActorBaseStat, T>
type NatureStats<T> = Record<Nature, T>

const actorFocuses = [
  'none',
  'aggressive',
  'relentless',
  'reckless',
  'heavy',
  'patient',
  'hardened',
  'tough',
  'steadfast',
  'intelligent',
  'volatile',
  'intense',
  'calculated',
  'composed',
  'mindful',
  'reserved',
  'stoic',
  'agile',
  'hasty',
  'impulsive',
  'alert',
] as const

type ActorFocus = (typeof actorFocuses)[number]
type ActorFocusDetail = {
  up: ActorNatureStat | null
  down: ActorNatureStat | null
}

const ACTOR_FOCUS_DETAILS: Record<ActorFocus, ActorFocusDetail> = {
  none: { up: null, down: null },
  aggressive: { up: 'attack', down: 'defense' },
  relentless: { up: 'attack', down: 'chakra_attack' },
  reckless: { up: 'attack', down: 'chakra_defense' },
  heavy: { up: 'attack', down: 'speed' },
  patient: { up: 'defense', down: 'attack' },
  hardened: { up: 'defense', down: 'chakra_attack' },
  tough: { up: 'defense', down: 'chakra_defense' },
  steadfast: { up: 'defense', down: 'speed' },
  intelligent: { up: 'chakra_attack', down: 'attack' },
  volatile: { up: 'chakra_attack', down: 'defense' },
  intense: { up: 'chakra_attack', down: 'chakra_defense' },
  calculated: { up: 'chakra_attack', down: 'speed' },
  composed: { up: 'chakra_defense', down: 'attack' },
  mindful: { up: 'chakra_defense', down: 'defense' },
  reserved: { up: 'chakra_defense', down: 'chakra_attack' },
  stoic: { up: 'chakra_defense', down: 'speed' },
  agile: { up: 'speed', down: 'attack' },
  hasty: { up: 'speed', down: 'defense' },
  impulsive: { up: 'speed', down: 'chakra_attack' },
  alert: { up: 'speed', down: 'chakra_defense' },
}

type ActorDef = {
  actor_ID: string
  sprite_url: string
  name: string
  clan: string
  affiliations: Array<string>
  restricted: boolean
  stats: ActorStats<number>
  natures: Partial<Record<NatureSet, Array<Nature>>>
  nature_damage: NatureStats<number>
  nature_resistance: NatureStats<number>
  default_modifiers: Array<Modifier>
  abilities: Array<Modifier>
  action_count: number
  action_IDs: Array<string>
}

type ActorState = {
  alive: boolean
  damage: number
  position_ID: string
  last_used_action_ID: string | null
  action_locked: boolean
  switch_locked: boolean
  seen: boolean
  stunned: boolean
  stamina_damage: number
  statused: boolean
  burned: boolean
  paralyzed: boolean
  poisoned: boolean
  sleeping: boolean
}

type Actor = ActorDef &
  ActorState & {
    ID: string
    player_ID: string
    enabled: boolean
    level: number
    experience: number
    focus: ActorFocus
    base_stats: ActorStats<number>
    staged_stats: ActorStats<number>
    unmodified_stats: ActorStats<number>
    aux_stats: ActorStats<number>
    applied_modifiers: Record<string, number>
    actions: Array<Action>
    resolved_nature_resistance: NatureStats<number>
    resolved_nature_damage: NatureStats<number>
    summon?: Actor & { proxy: boolean }
    ability: Modifier | null
    item: Modifier | null
  }

function checkActorStat(actor: Actor, key: ActorBaseStat) {
  const stat = actor.stats[key]
  const pre = actor.unmodified_stats[key]
  return stat === pre ? 0 : stat > pre ? 1 : -1
}

function getTotalBaseStats(actor: ActorDef) {
  const stats: ActorDef['stats'] = {
    ...actor.stats,
    accuracy: 0,
    evasion: 0,
  }
  return Math.floor(
    Object.values(stats).reduce((p, c) => p + c, 0) - stats.stamina
  )
}

function getVitals(actor: Actor) {
  const maxHp = Math.max(1, actor.stats.hp)
  const maxStamina = Math.max(1, actor.stats.stamina)
  const currentHp = Math.max(0, actor.stats.hp - actor.damage)
  const currentStamina = Math.max(0, actor.stats.stamina - actor.stamina_damage)

  const hpRatio = clamp(currentHp / maxHp)
  const staminaRatio = clamp(currentStamina / maxStamina)

  return {
    hp: {
      max: maxHp,
      current: currentHp,
      ratio: hpRatio,
    },
    stamina: {
      max: maxStamina,
      current: currentStamina,
      ratio: staminaRatio,
    },
  }
}

export type {
  ActorDef,
  Actor,
  ActorFocus,
  ActorAttackStat,
  ActorDefenseStat,
  ActorBaseStat,
  ActorStats,
  ActorNatureStat,
  NatureStats,
}
export {
  checkActorStat,
  getTotalBaseStats,
  actorFocuses,
  ACTOR_FOCUS_DETAILS,
  getVitals,
}
