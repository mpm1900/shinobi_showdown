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
    log_success: string
    log_failure: string
    target_type: 'target-actor-id' | 'target-position-type'
    summon: boolean
    switch: boolean
    struggle: boolean
  }
  state: {
    locked: boolean
    cooldown: number | null
    disabled: boolean
  }
  priority: number
}

type ActionTransaction = {
  ID: string
  context: Context
  mutation: Action
}

export type { Action, ActionTransaction }
