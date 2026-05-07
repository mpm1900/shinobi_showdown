import type { Actor } from '#/lib/game/actor'
import { cn } from '#/lib/utils'
import { CircleQuestionMark, X } from 'lucide-react'
import type { ComponentProps } from 'react'

function clamp(value: number) {
  return Math.max(0.02, Math.min(1, value))
}

function MiniHealthBar({
  className,
  actor,
  ...props
}: React.ComponentProps<'div'> & { actor: Actor }) {
  const maxHp = Math.max(1, actor.stats.hp)
  const hpCurrent = Math.max(0, actor.stats.hp - actor.damage)
  const hpRatio = clamp(hpCurrent / maxHp)

  return (
    <>
      <div
        {...props}
        className={cn(
          'absolute -bottom-px h-1.5 left-0 right-0 z-20 bg-black/80',
          className
        )}
      />
      <div
        {...props}
        className={cn(
          'absolute -bottom-px h-1.5 left-0 z-20 bg-red-400 border border-black',
          className
        )}
        style={{
          width: `${hpRatio * 100}%`,
        }}
      />
    </>
  )
}

function ActorThumbnail({
  actor,
  active,
  hidden,
  index,
  size = 64,
  className,
  showRing,
  showHealthBar,
  imgClass,
  ...props
}: ComponentProps<'div'> & {
  actor: Actor
  active?: boolean
  hidden?: boolean
  index?: number
  size?: number
  showRing?: boolean
  showHealthBar?: boolean
  imgClass?: string
}) {
  active = active ?? !!actor.position_ID
  const alive = actor.alive
  return (
    <div
      key={actor.ID}
      style={{
        height: size,
        width: size,
      }}
      className={cn(
        'overflow-hidden bg-stone-950 p-1 border rounded relative',
        {
          'bg-transparent border-transparent': index === undefined,
          'bg-foreground': active && index !== undefined,
          'ring ring-black': showRing,
        }
      )}
      {...props}
    >
      <div className={cn('h-full w-full select-none', className)}>
        {!alive && <div className="bg-red-950 absolute inset-0"></div>}
        {!hidden && (
          <img
            src={actor.sprite_url}
            draggable={false}
            className={cn(
              'h-full w-full object-cover absolute inset-0 z-10',
              imgClass
            )}
            style={{
              imageRendering: 'pixelated',
            }}
            width={size}
            height={size}
          />
        )}

        {!hidden && index !== undefined && (
          <div
            className={cn(
              'absolute -top-6 font-black text-7xl z-0 text-center text-foreground',
              {
                'text-background!': active,
              }
            )}
          >
            {index + 1}
          </div>
        )}
        {hidden && (
          <div className="grid place-items-center h-full w-full">
            <CircleQuestionMark className="size-12 text-muted-foreground/30" />
          </div>
        )}
        {!alive && (
          <div className="grid place-items-center absolute inset-0 z-20 ">
            <X className="size-12 text-destructive" />
          </div>
        )}
        {alive && !hidden && showHealthBar && <MiniHealthBar actor={actor} />}
      </div>
    </div>
  )
}

export { ActorThumbnail, MiniHealthBar }
