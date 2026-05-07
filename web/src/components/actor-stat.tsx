import {
  checkActorStat,
  type Actor,
  type ActorBaseStat,
  type ActorDef,
} from '#/lib/game/actor'
import type { Nature } from '#/lib/game/nature'
import { cn } from '#/lib/utils'

function getBaseStatColorClass(value: number) {
  if (value < 45) return 'text-rose-600'
  if (value < 55) return 'text-red-500'
  if (value < 65) return 'text-orange-600'
  if (value < 75) return 'text-orange-500'
  if (value < 85) return 'text-amber-500'
  if (value < 95) return 'text-yellow-500'
  if (value <= 100) return 'text-yellow-300'
  if (value < 110) return 'text-lime-300'
  if (value < 120) return 'text-lime-500'
  if (value < 130) return 'text-green-400'
  if (value < 140) return 'text-emerald-500'
  if (value < 150) return 'text-teal-500'
  if (value < 165) return 'text-teal-300'
  if (value < 180) return 'text-cyan-400'
  if (value < 200) return 'text-indigo-200'
  return 'text-fuchsia-300'
}

function ActorStat({
  stat,
  actor,
  className,
  showBase = true,
  ...props
}: React.ComponentProps<'span'> & {
  actor: Actor
  stat: ActorBaseStat
  showBase?: boolean
}) {
  const check = checkActorStat(actor, stat)
  return (
    <span
      data-role="actor-stat"
      className={cn('font-bold', className)}
      {...props}
    >
      <span
        className={cn('whitespace-nowrap', {
          'text-green-400': check === 1,
          'text-red-300': check === -1,
        })}
      >
        {actor.stats[stat]}
      </span>
      {showBase && (
        <span className="text-muted-foreground/60">
          {' '}
          ({actor.base_stats[stat]})
        </span>
      )}
    </span>
  )
}

function ActorStatBase({
  stat,
  actor,
  ...props
}: React.ComponentProps<'span'> & {
  actor: ActorDef
  stat: ActorBaseStat
}) {
  return (
    <span data-role="actor-stat" {...props}>
      <span className={getBaseStatColorClass(actor.stats[stat])}>
        {actor.stats[stat]}
      </span>
    </span>
  )
}

function NatureDamageStat({
  actor,
  nature,
  ...props
}: React.ComponentProps<'span'> & {
  actor: Actor
  nature: Nature
}) {
  const value = actor.resolved_nature_damage[nature]
  return (
    <span data-role="nature-stat" {...props}>
      <span
        className={cn('whitespace-nowrap', {
          'text-green-400': value > 1,
          'text-red-300': value < 1,
        })}
      >
        {100 * value}%
      </span>
    </span>
  )
}

function NatureResistanceStat({
  actor,
  nature,
  ...props
}: React.ComponentProps<'span'> & {
  actor: Actor
  nature: Nature
}) {
  const value = actor.resolved_nature_resistance[nature]
  return (
    <span data-role="nature-stat" {...props}>
      <span
        className={cn('whitespace-nowrap', {
          'text-green-400': value > 1 || value < 0,
          'text-red-300': value < 1 && value > 0,
        })}
      >
        {Math.floor(100 * value)}%
      </span>
    </span>
  )
}

export { ActorStat, ActorStatBase, NatureDamageStat, NatureResistanceStat }
