import { AppHeader } from '#/components/app-header'
import { BattleActions } from '#/components/battle-actions'
import { BattleContextController } from '#/components/battle-context-controller'
import { BattleTargetingVfx } from '#/components/battle-targeting-vfx'
import { BattleWeather } from '#/components/battle-weather'
import { GameLog } from '#/components/game-log'
import { PlayerPositions } from '#/components/player-positions'
import { PlayerThumbnails } from '#/components/player-thumbnails'
import { PromptController } from '#/components/prompt-controller'
import { RunningContext } from '#/components/running-context'
import { refreshVanta } from '#/components/vanta-background'
import { useActiveActor } from '#/hooks/use-active-actor'
import { battleContext, setContextSource } from '#/lib/stores/battle-context'
import { clientsStore } from '#/lib/stores/clients'
import { gameStore } from '#/lib/stores/game'
import { ClientOnly, createFileRoute, redirect } from '@tanstack/react-router'
import { useStore } from '@tanstack/react-store'
import { useEffect, useRef } from 'react'

export const Route = createFileRoute('/battle')({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: '/login' })
    }
  },
  component: RouteComponent,
})

function RouteComponent() {
  const game = useStore(gameStore, (g) => g)
  const client = useStore(clientsStore, (c) => c.me)
  const source_actor_ID = useStore(battleContext, (c) => c.source_actor_ID)
  const actor = useActiveActor()
  const players = game.players.filter((p) => p.ID === client?.ID)
  const enemies = game.players.filter((p) => p.ID !== client?.ID)
  const vfxContainerRef = useRef<HTMLDivElement>(null)
  const runningContext =
    game.status === 'running'
      ? (game.active_transaction?.context ?? null)
      : null

  useEffect(() => {
    setTimeout(() => refreshVanta(), 100)
  }, [])

  return (
    <ClientOnly>
      <PromptController />
      <BattleContextController />
      <main className="flex flex-col h-dvh overflow-hidden">
        <div className="hidden xl:block">
          <AppHeader />
        </div>
        <div
          ref={vfxContainerRef}
          className="flex flex-col flex-1 relative overflow-hidden"
        >
          <BattleTargetingVfx
            containerRef={vfxContainerRef}
            runningContext={runningContext}
          />
          <div>
            <div className="fixed top-2 xl:top-17 px-4 flex flex-col items-end right-4 z-30">
              <div>
                {enemies.map((player) => (
                  <PlayerPositions key={player.ID} flip player_ID={player.ID} />
                ))}
              </div>
              <GameLog />
            </div>
            <div className="fixed top-12 px-4 py-4 left-0 z-10">
              {enemies.map((player) => (
                <PlayerThumbnails key={player.ID} player_ID={player.ID} />
              ))}
              <BattleWeather />
            </div>
          </div>
          {game.turn.phase !== 'init' && (
            <div className="flex-1 grid place-items-center overflow-hidden relative">
              {actor && <BattleActions actor={actor} />}
              {game.status === 'running' &&
                game.active_transaction?.context && (
                  <RunningContext context={game.active_transaction?.context} />
                )}
            </div>
          )}
          <div className="fixed bottom-4 xl:bottom-8 left-8 flex z-40">
            {players.map((player) => (
              <PlayerPositions
                key={player.ID}
                flip={false}
                player_ID={player.ID}
                selected={game.status === 'idle' ? (source_actor_ID ?? '') : ''}
                onSelectedChange={(actor_ID) =>
                  setContextSource(actor_ID, game)
                }
              />
            ))}
          </div>
          <div className="fixed bottom-0 px-4 py-4 flex flex-col items-end right-0 z-30">
            {players.map((player) => (
              <PlayerThumbnails key={player.ID} player_ID={player.ID} />
            ))}
          </div>
        </div>
      </main>
    </ClientOnly>
  )
}
