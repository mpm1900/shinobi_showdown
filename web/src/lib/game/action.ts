import type { ActorNatureStat } from './actor'
import type { Context } from './context'
import type { NatureSet } from './nature'

type Action = {
  ID: string
  config: {
    name: string
    jutsu: string
    nature?: NatureSet
    cost?: number
    cooldown?: number
    accuracy?: number
    stat?: ActorNatureStat
    power?: number
    description: string
  }
  summon: boolean
  locked: boolean
  cooldown: number | null
  disabled: boolean
  priority: number
  target_type: 'target-actor-id' | 'target-position-type'
  meta: {
    switch: boolean
    struggle: boolean
  }
}

type ActionTransaction = {
  ID: string
  context: Context
  mutation: Action
}

export type { Action, ActionTransaction }
