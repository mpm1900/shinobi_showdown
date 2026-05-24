import { gameStore } from '#/lib/stores/game'
import { cn } from '#/lib/utils'
import { useStore } from '@tanstack/react-store'
import { ActorThumbnail } from './actor-thumbnail'

function LobbyThumbnails({
  className,
  player_ID,
  enabled = [],
  onEnabledChange,
}: {
  className?: string
  player_ID: string
  enabled?: string[]
  onEnabledChange?: (enabled: string[]) => void
}) {
  const actors = useStore(gameStore, (g) => g.actors)

  function toggleActor(ID: string) {
    if (enabled.length >= 4 && !enabled.includes(ID)) {
      return
    }
    if (enabled.includes(ID)) {
      onEnabledChange?.(enabled.filter((id) => id !== ID))
    } else {
      onEnabledChange?.([...enabled, ID])
    }
  }

  return (
    <div className={cn('flex gap-2', className)}>
      {actors
        .filter((a) => a.player_ID === player_ID)
        .map((a, i) => (
          <ActorThumbnail
            key={a.ID}
            actor={a}
            active={enabled.includes(a.ID)}
            onClick={() => toggleActor(a.ID)}
            index={i}
            showHealthBar
            showRing
            hidden={false}
            className={cn({
              'cursor-pointer': enabled.length < 4 || enabled.includes(a.ID),
              'cursor-not-allowed':
                enabled.length >= 4 && !enabled.includes(a.ID),
              'opacity-40': !enabled.includes(a.ID),
            })}
          />
        ))}
    </div>
  )
}

export { LobbyThumbnails }
