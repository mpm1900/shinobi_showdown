import { cn } from '#/lib/utils'
import { useStore } from '@tanstack/react-store'
import { ActorThumbnail } from './actor-thumbnail'
import { gameStore } from '#/lib/stores/game'
import { clientsStore } from '#/lib/stores/clients'
import { ActorTooltip } from './actor-tooltip'

function PlayerThumbnails({ player_ID }: { player_ID: string }) {
  const client = useStore(clientsStore, (s) => s.me)
  const actors = useStore(gameStore, (g) => g.actors)
  const is_player = client?.ID == player_ID
  return (
    <div className="flex gap-2">
      {actors
        .filter((a) => a.player_ID === player_ID)
        .map((a, i) => (
          <ActorTooltip key={a.ID} actor={a} disabled={!is_player}>
            <ActorThumbnail
              actor={a}
              index={i}
              showHealthBar={is_player}
              showRing
              hidden={!is_player && !a.seen && !a.position_ID}
              className={cn({
                'opacity-40': !a.position_ID,
              })}
            />
          </ActorTooltip>
        ))}
    </div>
  )
}

export { PlayerThumbnails }
