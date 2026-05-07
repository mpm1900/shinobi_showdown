import type { Context } from './context'

type Modifier = {
  ID: string
  group_ID: string
  name: string
  description: string
  parent_description: string
  duration: number | null
  icon?: string
  show: boolean
}

type ModifierTransaction = {
  ID: string
  context: Context
  mutation: Modifier
}

export type { Modifier, ModifierTransaction }
