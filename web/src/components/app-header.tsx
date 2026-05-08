import { useStore } from '@tanstack/react-store'
import { InstanceCombobox } from './instance-combobox'
import { sendContextMessage, socketStore } from '#/lib/stores/socket'
import { clientsStore } from '#/lib/stores/clients'
import { gameStore } from '#/lib/stores/game'
import {
  ChevronRight,
  Loader,
  LogOut,
  Signal,
  TriangleAlert,
  Unplug,
} from 'lucide-react'
import { Tabs, TabsList, TabsTrigger } from './ui/tabs'
import { Link, useRouterState } from '@tanstack/react-router'
import { NULL_CONTEXT } from '#/lib/game/context'
import { Button } from './ui/button'
import { GiSharpShuriken } from 'react-icons/gi'
import { useLogout } from '#/lib/mutations/logout'
import { useUser } from '#/lib/queries/auth'
import { Tooltip, TooltipContent, TooltipTrigger } from './ui/tooltip'
import { GiWhirlpoolShuriken } from 'react-icons/gi'
import {
  TbHexagonNumber1Filled,
  TbHexagonNumber2Filled,
  TbHexagonNumber3Filled,
} from 'react-icons/tb'
import { connect } from '#/lib/socket/connect'

function getActiveTable(pathname: string) {
  switch (pathname) {
    case '/battle':
      return 'battle'
    case '/team-builder':
      return 'team-builder'
    case '/lobby':
      return 'lobby'
    default:
      return 'debug'
  }
}

function AppHeader() {
  const { data: user } = useUser()
  const logout = useLogout()
  const instanceID = useStore(socketStore, (s) => s.instanceID)
  const status = useStore(socketStore, (s) => s.status)
  const client = useStore(clientsStore, (c) => c.me)
  const game_status = useStore(gameStore, (g) => g.status)
  const game_phase = useStore(gameStore, (g) => g.turn.phase)
  const actions = useStore(gameStore, (g) => g.actions)
  const pathname = useRouterState({
    select: (state) => state.location.pathname,
  })
  const activeTab = getActiveTable(pathname)
  return (
    <header className="fixed flex justify-between p-1 z-20 bg-stone-950 border-b border-stone-800 ring ring-black w-full">
      <div className="flex items-center gap-2">
        <Link to="/" className="pl-2">
          <GiSharpShuriken />
        </Link>

        {user && (
          <InstanceCombobox
            icon={
              <>
                {status === 'idle' && <Unplug />}
                {(status === 'connecting' || status === 'reconnecting') && (
                  <Loader className="animate-spin" />
                )}
                {status === 'open' && <Signal />}
                {(status === 'closed' || status === 'error') && (
                  <TriangleAlert className="text-destructive" />
                )}
              </>
            }
            value={instanceID}
            onValueChange={(instanceID) => connect(instanceID)}
          />
        )}
        <Tabs value={activeTab}>
          <TabsList className="ring-0">
            <TabsTrigger value="team-builder" asChild>
              <Link to="/team-builder">
                <TbHexagonNumber1Filled />
                Team Builder
              </Link>
            </TabsTrigger>
            <ChevronRight />
            <TabsTrigger value="lobby" asChild>
              <Link to="/lobby">
                <TbHexagonNumber2Filled />
                Lobby
              </Link>
            </TabsTrigger>
            <ChevronRight />
            <TabsTrigger value="battle" asChild>
              <Link to="/battle">
                <TbHexagonNumber3Filled />
                Battle
              </Link>
            </TabsTrigger>
          </TabsList>
        </Tabs>

        {client && (
          <div className="flex gap-2">
            {game_phase !== 'init' && (
              <Button
                disabled={game_status === 'running'}
                onClick={() => {
                  sendContextMessage({
                    type: 'run-game-actions',
                    client_ID: client.ID,
                    context: NULL_CONTEXT,
                  })
                }}
              >
                Run
              </Button>
            )}
          </div>
        )}
        <div className="flex items-center">
          {game_status === 'running' && (
            <GiWhirlpoolShuriken className="animate-spin" />
          )}
          {game_status === 'waiting' && (
            <GiWhirlpoolShuriken className="animate-spin" />
          )}
        </div>
      </div>
      <div className="flex items-center gap-4 px-2">
        <div className="font-mono text-sm flex items-center">
          {user && (
            <div className="flex items-center gap-2">
              <Tooltip>
                <TooltipTrigger>
                  <span>{user.email}</span>
                </TooltipTrigger>
                <TooltipContent>{user.id}</TooltipContent>
              </Tooltip>
              <Button
                variant="ghost"
                size="icon"
                className="size-8"
                onClick={() => logout.mutate()}
                title="Logout"
              >
                <LogOut className="size-4" />
              </Button>
            </div>
          )}
        </div>
      </div>
    </header>
  )
}

export { AppHeader }
