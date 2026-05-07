import z from 'zod'
import { actorFocuses } from '../game/actor'

const ActorConfigSchema = z
  .object({
    ability_ID: z.string().nullable(),
    action_IDs: z.array(z.string()).min(1),
    focus: z.enum(actorFocuses),
    item_ID: z.string().nullable(),
    aux_stats: z.object({
      hp: z.number().min(0, 'negative').max(31, 'too big'),
      attack: z.number().min(0, 'negative').max(31, 'too big'),
      stamina: z.number().min(0, 'negative').max(31, 'too big'),
      defense: z.number().min(0, 'negative').max(31, 'too big'),
      speed: z.number().min(0, 'negative').max(31, 'too big'),
      chakra_attack: z.number().min(0, 'negative').max(31, 'too big'),
      chakra_defense: z.number().min(0, 'negative').max(31, 'too big'),
    }),
  })
  .superRefine((config, ctx) => {
    const total = Object.values(config.aux_stats).reduce(
      (sum, value) => sum + value,
      0
    )
    if (total > 66) {
      ctx.addIssue({
        code: 'custom',
        message: 'Total aux stats cannot exceed 66',
      })
    }
  })

const TeamActorSchema = z.object({
  actor_ID: z.string(),
  config: ActorConfigSchema,
})

const TeamSchema = z
  .object({
    name: z.string().min(3),
    actors: z.array(TeamActorSchema).min(1),
    selected_index: z.number(),
  })
  .superRefine((team, ctx) => {
    const hasActiveActor = team.actors.some((a) => a !== null)
    if (!hasActiveActor) {
      ctx.addIssue({
        code: 'custom',
        path: ['actors'],
        message: 'Must have an active actor',
      })
    }
  })

type ActorConfig = z.output<typeof ActorConfigSchema>
type TeamConfig = z.output<typeof TeamSchema>
type TeamActor = z.output<typeof TeamActorSchema>

export { ActorConfigSchema, TeamActorSchema, TeamSchema }
export type { ActorConfig, TeamConfig, TeamActor }
