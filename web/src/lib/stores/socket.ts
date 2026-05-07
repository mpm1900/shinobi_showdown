import { Store } from '@tanstack/store'
import type { ActorConfig } from './config'
import type {
  SocketMessageSubscriber,
  SocketRequest,
  SocketResponse,
} from '../socket/request'

type SocketStatus =
  | 'idle'
  | 'connecting'
  | 'open'
  | 'closing'
  | 'closed'
  | 'error'
  | 'reconnecting'

type SocketState = {
  instanceID: string | null
  socket: WebSocket | null
  status: SocketStatus
  reconnectCount: number
  isManualDisconnect: boolean
}

const socketStore = new Store<SocketState>({
  instanceID: null,
  socket: null,
  status: 'idle',
  reconnectCount: 0,
  isManualDisconnect: false,
})

function sendSocketMessage(
  payload: string | ArrayBufferLike | Blob | ArrayBufferView
): boolean {
  const socket = socketStore.state.socket
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    console.warn('Attempted to send message while socket is not open')
    return false
  }

  socket.send(payload)
  return true
}

function sendContextMessage(request: SocketRequest) {
  return sendSocketMessage(JSON.stringify(request))
}

export type { SocketResponse, SocketMessageSubscriber, ActorConfig }
export { socketStore, sendSocketMessage, sendContextMessage }
