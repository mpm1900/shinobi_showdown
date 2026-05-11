import { SHINOBI_ICONS } from '#/data/icons'
import type { Action } from '#/lib/game/action'
import { type Actor } from '#/lib/game/actor'
import type { Context } from '#/lib/game/context'
import { addHoverTarget, removeHoverTarget } from '#/lib/stores/battle-context'
import { cn, keys } from '#/lib/utils'
import { useStore } from '@tanstack/react-store'
import { MiniHealthBar } from './actor-thumbnail'
import { Button } from './ui/button'
import { gameStore } from '#/lib/stores/game'
import { getEffectiveness } from '#/lib/game/nature'
import { memo, useMemo } from 'react'

const TargetButton = memo(function TargetButton({
  actor,
  context,
  contextValid,
  enabled,
  loading,
  onContextChange,
  targetType,
}: {
  actor: Actor
  context: Context
  contextValid: boolean
  enabled: boolean
  loading: boolean
  onContextChange: (context: Context) => void
  targetType: Action['target_type']
}) {
  const source = useStore(gameStore, (g) =>
    g.actors.find((a) => a.ID === context.source_actor_ID)
  )
  const action = source?.actions?.find((a) => a.ID === context.action_ID)
  const effectiveness = useMemo(() => {
    return action?.config.nature
      ? getEffectiveness(action.config.nature, keys(actor.natures))
      : null
  }, [action?.config.nature, actor.natures])

  const includes =
    targetType === 'target-actor-id'
      ? context.target_actor_IDs?.includes(actor.ID)
      : context.target_position_IDs?.includes(actor.position_ID)

  return (
    <Button
      className="relative flex-col h-auto p-2 px-3 min-w-56 w-auto overflow-hidden shadow-[4px_4px_8px_rgba(0,0,0,0.5)]"
      disabled={loading || (contextValid && !includes) || !enabled}
      variant={
        includes
          ? 'default'
          : context.source_player_ID === actor.player_ID
            ? 'player_target'
            : 'enemy_target'
      }
      onMouseEnter={() => {
        addHoverTarget(actor.ID)
      }}
      onMouseLeave={() => {
        removeHoverTarget(actor.ID)
      }}
      onClick={() => {
        if (targetType === 'target-actor-id') {
          onContextChange({
            ...context,
            target_actor_IDs: includes
              ? (context.target_actor_IDs?.filter((id) => id !== actor.ID) ??
                null)
              : [...(context.target_actor_IDs ?? []), actor.ID],
          })
        }

        if (targetType === 'target-position-type') {
          onContextChange({
            ...context,
            target_position_IDs: includes
              ? (context.target_position_IDs?.filter(
                (id) => id !== actor.position_ID
              ) ?? null)
              : [...(context.target_position_IDs ?? []), actor.position_ID],
          })
        }
      }}
    >
      <div
        className={cn(
          'flex items-end w-full justify-between gap-4 relative z-10',
          !includes && 'text-shadow-[1px_1px_0px_#000000]'
        )}
      >
        <div className="truncate font-bold">{actor.name}</div>
        {action?.config.power && (
          <div className="text-xs">
            x<span className="font-black">{effectiveness?.toFixed(2)}</span>
          </div>
        )}
      </div>

      <div className="relative w-full">
        <MiniHealthBar actor={actor} className="left-0 right-0" />
      </div>
      <img
        src={actor.sprite_url}
        draggable={false}
        className={cn('absolute left-0 bottom-0 opacity-40')}
        style={{
          imageRendering: 'pixelated',
        }}
        width={64}
        height={64}
      />
    </Button>
  )
})

export { TargetButton }
