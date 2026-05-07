import type { Actor } from '#/lib/game/actor'
import { natureIndexes } from '#/lib/game/nature'
import { cn, keys } from '#/lib/utils'
import { cva, type VariantProps } from 'class-variance-authority'
import { HealthBar } from './health-bar'
import { NatureBadge } from './nature-badge'
import { Badge } from './ui/badge'
import { X } from 'lucide-react'
import { sendContextMessage } from '#/lib/stores/socket'
import { NULL_CONTEXT } from '#/lib/game/context'
import { ActorThumbnail } from './actor-thumbnail'
import { StageBadge } from './stage-badge'
import { ActorModifiers } from './actor-modifiers'
import { ActorStatus } from './actor-status'
import { useStore } from '@tanstack/react-store'
import { gameStore } from '#/lib/stores/game'
import { useDebounce } from 'use-debounce'

const frameVariants = cva('border border-transparent', {
  variants: {
    variant: {
      default: 'bg-stone-900 text-stone-200 border-stone-700',
      selected:
        'bg-stone-300 text-stone-900 border-stone-100 text-shadow-none!',
      targeted: 'bg-red-900 text-stone-200 border-red-500',
    },
  },
  defaultVariants: {
    variant: 'default',
  },
})

function parseSummon(summon: Actor['summon'], parent: Actor): Actor['summon'] {
  if (!summon) return summon
  if (summon.proxy) {
    return {
      ...summon,
      ...parent,
      ID: summon.ID,
      stats: {
        ...parent.stats,
        hp: summon.stats.hp,
      },
      damage: summon.damage,
      sprite_url: summon.sprite_url,
      applied_modifiers: {},
      summon: undefined,
    }
  }

  return summon
}

type ActorCardProps = {
  actor: Actor | undefined
  client_ID?: string
  selected?: boolean
  targeted?: boolean
  source?: boolean
  summon?: boolean
  summonClass?: string
}

type Variant = VariantProps<typeof frameVariants>['variant']

function getVariant(props: Partial<ActorCardProps>): Variant {
  if (props.targeted) return 'targeted'
  if (props.selected) return 'selected'
  return 'default'
}

function ActorCard({
  actor,
  client_ID,
  selected,
  source,
  targeted,
  className,
  summon,
  summonClass,
  ...rest
}: React.ComponentProps<'div'> & ActorCardProps) {
  const status = useStore(gameStore, (s) => s.status)
  const actions = useStore(gameStore, (s) => s.actions)
  const modifiers = useStore(gameStore, (g) =>
    (g.modifiers ?? [])
      .map((m) => m.mutation)
      .concat(g.actors.filter((a) => a.ability).map((a) => a.ability!))
      .concat(g.actors.filter((a) => a.item).map((a) => a.item!))
      .concat(
        g.actors
          .filter((a) => a.default_modifiers)
          .flatMap((a) => a.default_modifiers)
      )
  )
  const has_queued_action = useStore(
    gameStore,
    (s) => s.queued_actions[actor?.ID ?? '']
  )
  const action_tx = actions.find((t) => t.context.source_actor_ID === actor?.ID)
  const is_player = actor?.player_ID === client_ID
  const v = getVariant({ selected, targeted })
  const [variant] = useDebounce(v, 200)

  return (
    <div
      className={cn(
        'relative flex flex-col',
        summon && 'pointer-events-none',
        className
      )}
    >
      {actor?.summon?.proxy && (
        <ActorCard
          summon
          actor={parseSummon(actor.summon, actor)}
          client_ID={client_ID}
          selected={selected}
          targeted={targeted}
          source={source}
          className={cn('absolute top-0 left-0 z-20', summonClass)}
        />
      )}
      {actor && <ActorModifiers actor={actor} modifiers={modifiers} />}
      <div
        className={cn(
          is_player && 'cursor-pointer',
          'group/item flex items-center rounded-md text-sm transition-colors duration-100 outline-none focus-visible:border-ring focus-visible:ring-[3px] focus-visible:ring-ring/50 [a]:transition-colors [a]:hover:bg-accent/50 px-2 pt-0 pb-3 w-108'
        )}
        {...rest}
      >
        <ActorAvatar actor={actor} variant={variant} />
        <div className="flex flex-1 flex-col relative gap-0">
          {actor && (
            <div className="flex justify-between items-end gap-4">
              <div
                className={cn(
                  'pl-2 pr-4 pb-2 pt-1.5 -mb-3 rounded-none rounded-tr-xl shadow-[2px_1px_0px_rgba(0,0,0,1)]',
                  frameVariants({
                    variant,
                  }),
                  'border-l-0'
                )}
              >
                <span
                  className={cn(
                    'font-semibold text-2xl text-nowrap',
                    !selected && 'text-shadow-[1px_1px_0px_#000000]',
                    'nanum-brush-script-regular select-none'
                  )}
                >
                  {actor.name}
                </span>
              </div>
              {!action_tx || !is_player || status === 'running' ? (
                <ActorNatures actor={actor} />
              ) : (
                <Badge
                  className="mb-0.5 max-w-36"
                  onClick={(e) => {
                    e.preventDefault()
                    e.stopPropagation()
                    if (has_queued_action) return

                    sendContextMessage({
                      type: 'remove-action',
                      client_ID: client_ID!,
                      context: {
                        ...NULL_CONTEXT,
                        action_ID: action_tx.ID,
                      },
                    })
                  }}
                >
                  <span className="truncate">
                    {action_tx.mutation.config.name}{' '}
                  </span>
                  {!has_queued_action && <X />}
                </Badge>
              )}
            </div>
          )}
          <div className="space-y-2">
            {actor && (
              <HealthBar
                actor={actor}
                selected={selected}
                percentage={!is_player}
              />
            )}
          </div>
          <ActorStages actor={actor} />
        </div>
      </div>
    </div>
  )
}

function ActorStages({ actor }: { actor: Actor | undefined }) {
  if (!actor) return null

  return (
    <div className="absolute -bottom-2.5 flex gap-1 px-2">
      {Object.entries(actor.staged_stats)
        //.filter(([key]) => key !== 'evasion' && key !== 'accuracy')
        .map(([key, stage]) => (
          <StageBadge key={key} stage={stage as any} stat={key as any} />
        ))}
    </div>
  )
}

function ActorNatures({ actor }: { actor: Actor | undefined }) {
  if (!actor) return null
  return (
    <div className="flex gap-px pb-1">
      {keys(actor.natures)
        .sort((a, b) => natureIndexes[a] - natureIndexes[b])
        .map((nature) => (
          <NatureBadge key={nature} nature={nature} />
        ))}
    </div>
  )
}

function ActorAvatar({
  actor,
  variant,
}: {
  actor: Actor | undefined
  variant: Variant
}) {
  if (!actor) return null

  return (
    <div
      data-actor-avatar-anchor={actor.ID}
      className={cn(
        'relative shrink-0 px-1 pb-2 -mb-3.5 -mr-1 rounded-3xl rounded-tr-none shadow - [2px_1px_0px_rgba(0, 0, 0, 1)] supports - [font: -apple - system - body]: translate - y - px',
        frameVariants({
          variant: variant,
        }),
        'border-r-0'
      )}
    >
      <ActorThumbnail actor={actor} size={64} imgClass="rounded-bl-xl" />
      <ActorStatus actor={actor} />
    </div>
  )
}

export { ActorCard }
