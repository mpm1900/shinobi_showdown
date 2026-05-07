import z from 'zod'
import type { Game } from './game'
import type { Actor } from './actor'
import type { Action } from './action'

const ContextSchema = z.object({
  action_ID: z.string().nullable(),
  source_player_ID: z.string().nullable(),

  parent_actor_ID: z.string().nullable(),
  source_actor_ID: z.string().nullable(),

  target_actor_IDs: z.array(z.string()).nullable(),
  target_position_IDs: z.array(z.string()).nullable(),
})

type Context = z.output<typeof ContextSchema>

function contextToString(c: Context): string {
  return `${c.action_ID}.${c.parent_actor_ID}.${c.source_actor_ID}.${c.source_player_ID}.${c.target_actor_IDs?.join('+')}.${c.target_position_IDs?.join('+')}`
}

function getTargets(type: Action['target_type'] | undefined, game: Game, context: Context): Actor[] {
  const t_targets = game.actors.filter((a) =>
    context.target_actor_IDs?.includes(a.ID)
  )
  const p_targets = game.actors.filter((a) =>
    context.target_position_IDs?.includes(a.position_ID)
  )
  if (type === 'target-actor-id') {
    return t_targets
  }
  if (type === 'target-position-type') {
    return p_targets
  }
  return [...t_targets, ...p_targets]
}

const NULL_CONTEXT: Context = {
  action_ID: null,
  parent_actor_ID: null,
  source_actor_ID: null,
  source_player_ID: null,
  target_actor_IDs: [],
  target_position_IDs: [],
}

export { ContextSchema, contextToString, getTargets, NULL_CONTEXT }
export type { Context }
