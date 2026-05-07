import { useStore } from '@tanstack/react-store'
import {
  ClientOnly,
  createFileRoute,
  Link,
  redirect,
} from '@tanstack/react-router'
import { gameStore } from '#/lib/stores/game'
import { clientsStore } from '#/lib/stores/clients'
import { PromptController } from '#/components/prompt-controller'
import { AppHeader } from '#/components/app-header'
import { CardDescription, CardTitle } from '#/components/ui/card'
import { sendContextMessage } from '#/lib/stores/socket'
import { NULL_CONTEXT } from '#/lib/game/context'
import { Button } from '#/components/ui/button'
import { useEffect, useState } from 'react'
import { LobbyThumbnails } from '#/components/lobby-thumbnails'
import { LobbyActorsList } from '#/components/lobby-actors-list'
import { Check, Loader2, Swords } from 'lucide-react'
import { LobbyTeamSelect } from '#/components/lobby-team-select'

export const Route = createFileRoute('/lobby')({
  beforeLoad: ({ context }) => {
    if (!context.auth.user) {
      throw redirect({ to: '/login' })
    }
  },
  component: App,
})

function App() {
  const client = useStore(clientsStore, (c) => c.me)
  const game = useStore(gameStore, (g) => g)
  const players = game.players.filter((p) => p.ID === client?.ID)
  const player = players[0]
  const enemies = game.players.filter((p) => p.ID !== client?.ID)
  const ready = players.length > 0
  const unstarted = game.status !== 'running' && game.turn.count == 0
  const nav = Route.useNavigate()
  const pids = players.map((p) => p.ID)
  const default_enabled = game.actors.filter(
    (a) => pids.includes(a.player_ID) && a.enabled
  )
  const [enabled, setEnabled] = useState<string[]>(
    default_enabled.map((a) => a.ID)
  )

  useEffect(() => {
    if (game.status === 'running') {
      nav({ to: '/battle' })
    }
  }, [game.status])

  return (
    <ClientOnly>
      <PromptController />
      <main className="min-w-0 overflow-x-hidden flex-1 flex flex-col">
        <AppHeader />
        <div className="flex min-w-0 flex-1">
          <div className="min-w-0 space-y-2 flex-1 overflow-auto m-10 my-6">
            <CardTitle>Pre-Game Lobby</CardTitle>
            <CardDescription>Select 4-shinobi your line-up</CardDescription>
            <div className="flex gap-2 mb-4">
              {client &&
                ready &&
                (unstarted ? (
                  <Link to="/battle">
                    <Button
                      disabled={
                        players.some((p) => !p.ready) ||
                        enemies.some((e) => !e.ready)
                      }
                      onClick={() => {
                        sendContextMessage({
                          type: 'start-battle',
                          client_ID: client!.ID,
                          context: NULL_CONTEXT,
                        })
                      }}
                    >
                      <Swords /> Start Battle
                    </Button>
                  </Link>
                ) : (
                  <Button
                    onClick={() => {
                      sendContextMessage({
                        type: 'reset',
                        client_ID: client!.ID,
                        context: NULL_CONTEXT,
                      })
                    }}
                  >
                    Reset
                  </Button>
                ))}
              {game.turn.phase === 'init' && !player?.ready && (
                <Button
                  disabled={enabled.length !== 4}
                  onClick={() => {
                    sendContextMessage({
                      type: 'ready-team',
                      client_ID: client!.ID,
                      context: {
                        ...NULL_CONTEXT,
                        target_actor_IDs: enabled,
                      },
                    })
                  }}
                >
                  Ready Team ({enabled.length}/4)
                </Button>
              )}
              {game.turn.phase === 'init' && player?.ready && (
                <Button
                  variant="destructive"
                  onClick={() => {
                    sendContextMessage({
                      type: 'cancel-team',
                      client_ID: client!.ID,
                      context: NULL_CONTEXT,
                    })
                  }}
                >
                  Cancel
                </Button>
              )}
              <LobbyTeamSelect
                onValueChange={() => {
                  setEnabled([])
                }}
              />
            </div>
            <div className="grid grid-cols-2 gap-12">
              {player && (
                <div className="flex flex-col flex-1 gap-2">
                  <div className="flex items-end justify-start">
                    <h3 className="font-bold text-xl">
                      {player.user.email} (You)
                    </h3>
                  </div>
                  <div className="flex items-center justify-between">
                    <LobbyThumbnails
                      player_ID={player.ID}
                      enabled={enabled}
                      onEnabledChange={setEnabled}
                    />
                    <div className="px-4">
                      {player.ready ? (
                        <Check className="text-green-300" />
                      ) : (
                        <Loader2 className="animate-spin text-muted-foreground" />
                      )}
                    </div>
                  </div>

                  <LobbyActorsList
                    player={player}
                    enabled={enabled}
                    onEnabledChange={setEnabled}
                  />
                </div>
              )}
              {enemies.map((player) => (
                <div key={player.ID} className="flex flex-col flex-1 gap-2">
                  <div className="flex justify-end item-end text-right">
                    <h3 className="font-bold text-xl">{player.user.email}</h3>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="px-4 ">
                      {player.ready ? (
                        <Check className="text-green-300" />
                      ) : (
                        <Loader2 className="animate-spin text-muted-foreground" />
                      )}
                    </div>
                    <LobbyThumbnails
                      player_ID={player.ID}
                      className="justify-end"
                    />
                  </div>
                  <LobbyActorsList player={player} />
                </div>
              ))}
            </div>
          </div>
        </div>
      </main>
    </ClientOnly>
  )
}
