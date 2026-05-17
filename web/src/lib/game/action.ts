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
  }
  state: {
    locked: boolean
    cooldown: number | null
    disabled: boolean
  }
  priority: number
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
