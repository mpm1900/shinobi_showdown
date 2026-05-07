import { useStore } from '@tanstack/react-store'
import { LobbyActorDetails } from './lobby-actor-details'
import { gameStore } from '#/lib/stores/game'
import type { Player } from '#/lib/game/player'

function LobbyActorsList({
  player,
  enabled = [],
  onEnabledChange,
}: {
  player: Player
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
    <>
      {actors
        .filter((a) => a.player_ID === player.ID)
        .map((a) => (
          <LobbyActorDetails
            key={a.ID}
            actor={a}
            enabled={enabled.includes(a.ID)}
            onClick={() => toggleActor(a.ID)}
          />
        ))}
    </>
  )
}

export { LobbyActorsList }
