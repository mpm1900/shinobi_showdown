import type { Context } from '../game/context'
import type { Game } from '../game/game'
import type { Client } from '../stores/clients'
import type { ActorConfig, TeamConfig } from '../stores/config'

type SocketRequestType =
  | 'set-team'
  | 'ready-team'
  | 'cancel-team'
  | 'start-battle'
  | 'reset'
  | 'push-action'
  | 'remove-action'
  | 'run-game-actions'
  | 'resolve-prompt'
  | 'validate-context'
  | 'get-targets'
  | 'send-chat'

type SocketRequest = {
  type: SocketRequestType
  prompt_ID?: string
  client_ID: string
  context: Context
  actor_config?: Partial<ActorConfig>
  team_config?: TeamConfig
}

type ChatMessage = {
  client_ID: string
  text: string
  timestamp: string
  username: string
}

type SocketResponse = {
  type:
    | 'game'
    | 'clients'
    | 'join-success'
    | 'validate-context'
    | 'target-IDs'
    | 'new-chat'
  state: Game | null
  clients: Array<Client> | null
  context: Context | null
  valid: boolean | null
  chat_message: ChatMessage | null
}

type SocketMessageSubscriber = (
  event: MessageEvent,
  message: SocketResponse | null
) => void

export type {
  SocketMessageSubscriber,
  SocketRequest,
  SocketRequestType,
  SocketResponse,
}
