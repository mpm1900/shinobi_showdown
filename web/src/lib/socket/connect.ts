import {
  socketStore,
  type SocketMessageSubscriber,
  type SocketResponse,
} from '../stores/socket'
import {
  INITIAL_RECONNECT_DELAY,
  INSTANCE_ID_KEY,
  MAX_RECONNECT_DELAY,
} from './config'
import { messageReducer } from './message-reducer'

function getSocketUrl(instanceID: string): string {
  const envUrl = import.meta.env.VITE_BACKEND_URL
  if (envUrl) {
    const protocol = envUrl.startsWith('https') ? 'wss' : 'ws'
    const host = envUrl.replace(/^https?:\/\//, '')
    return `${protocol}://${host}/socket/${instanceID}/connect`
  }

  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const hostname = window.location.hostname
  const port = '3005'

  return `${protocol}://${hostname}:${port}/socket/${instanceID}/connect`
}

function clearSocketEventHandlers(socket: WebSocket) {
  socket.onopen = null
  socket.onclose = null
  socket.onerror = null
  socket.onmessage = null
}

function isCurrentSocket(socket: WebSocket): boolean {
  return socketStore.state.socket === socket
}

const messageSubscribers = new Set<SocketMessageSubscriber>()
let connectionAbortController: AbortController | null = null
let reconnectTimer: number | null = null

function attemptReconnect() {
  const { instanceID, reconnectCount, isManualDisconnect } = socketStore.state
  if (!instanceID || isManualDisconnect) return

  const delay = Math.min(
    INITIAL_RECONNECT_DELAY * Math.pow(2, reconnectCount),
    MAX_RECONNECT_DELAY
  )

  console.log(
    `Attempting to reconnect in ${delay}ms... (attempt ${reconnectCount + 1})`
  )

  socketStore.setState((s) => ({
    ...s,
    reconnectCount: s.reconnectCount + 1,
    status: 'reconnecting',
  }))

  reconnectTimer = window.setTimeout(() => {
    reconnectTimer = null
    connect(instanceID)
  }, delay)
}

function connect(instanceID: string, onOpen?: () => void) {
  const store = socketStore.get()
  if (!instanceID) return

  // Cancel any existing connection attempts
  if (connectionAbortController) {
    connectionAbortController.abort()
  }
  connectionAbortController = new AbortController()
  const signal = connectionAbortController.signal

  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }

  const previous = store.socket
  if (previous) {
    clearSocketEventHandlers(previous)
    if (
      previous.readyState === WebSocket.CONNECTING ||
      previous.readyState === WebSocket.OPEN
    ) {
      previous.close(1000, 'Switching instance')
    }
  }

  localStorage.setItem(INSTANCE_ID_KEY, instanceID)
  const url = getSocketUrl(instanceID)
  const socket = new WebSocket(url)

  socketStore.setState((s) => ({
    ...s,
    instanceID,
    socket,
    status: s.reconnectCount > 0 ? 'reconnecting' : 'connecting',
    isManualDisconnect: false,
  }))

  socket.onopen = () => {
    if (signal.aborted || !isCurrentSocket(socket)) {
      socket.close(1000, 'Aborted')
      return
    }
    console.log('WebSocket connection opened')
    socketStore.setState((s) => ({
      ...s,
      status: 'open',
    }))

    // Only reset reconnect count after a stable connection (e.g., 5 seconds)
    setTimeout(() => {
      if (
        !signal.aborted &&
        socketStore.state.socket === socket &&
        socketStore.state.status === 'open'
      ) {
        console.log('Connection stable, resetting reconnect count')
        socketStore.setState((s) => ({
          ...s,
          reconnectCount: 0,
        }))
      }
    }, 5000)

    onOpen?.()
  }

  socket.onmessage = (event) => {
    if (signal.aborted || !isCurrentSocket(socket)) return

    let message: SocketResponse | null = null

    if (typeof event.data === 'string') {
      try {
        message = JSON.parse(event.data) as SocketResponse
      } catch {
        console.error('non-JSON payload')
      }
    }

    messageReducer(message)

    for (const subscriber of messageSubscribers) {
      try {
        subscriber(event, message)
      } catch (error) {
        console.error('socket message subscriber error', error)
      }
    }
  }

  socket.onerror = (error) => {
    if (signal.aborted || !isCurrentSocket(socket)) return
    console.error('WebSocket error:', error)
    socketStore.setState((s) => ({
      ...s,
      status: 'error',
    }))
  }

  socket.onclose = (event) => {
    console.log(
      `WebSocket connection closed: code=${event.code}, reason=${event.reason}, wasClean=${event.wasClean}`
    )
    if (signal.aborted || !isCurrentSocket(socket)) return

    const { isManualDisconnect } = socketStore.state

    socketStore.setState((s) => ({
      ...s,
      socket: null,
      status: isManualDisconnect ? 'closed' : 'error',
    }))

    if (!isManualDisconnect) {
      attemptReconnect()
    }
  }
}

function disconnect(code = 1000, reason = 'Manual disconnect') {
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }

  const socket = socketStore.state.socket
  if (!socket) {
    socketStore.setState((s) => ({
      ...s,
      status: 'closed',
      instanceID: null,
      reconnectCount: 0,
      isManualDisconnect: true,
    }))
    if (typeof window !== 'undefined') {
      localStorage.removeItem(INSTANCE_ID_KEY)
    }
    return
  }

  socketStore.setState((s) => ({
    ...s,
    status: 'closing',
    isManualDisconnect: true,
  }))

  if (
    socket.readyState === WebSocket.CONNECTING ||
    socket.readyState === WebSocket.OPEN
  ) {
    socket.close(code, reason)
  } else {
    socketStore.setState((s) => ({
      ...s,
      socket: null,
      status: 'closed',
      instanceID: null,
      reconnectCount: 0,
    }))
  }

  localStorage.removeItem(INSTANCE_ID_KEY)
}

function subscribe(subscriber: SocketMessageSubscriber) {
  messageSubscribers.add(subscriber)
  return () => {
    messageSubscribers.delete(subscriber)
  }
}

function readSavedInstanceID(): string | null {
  if (typeof window === 'undefined') {
    return null
  }
  return localStorage.getItem(INSTANCE_ID_KEY)
}

if (typeof window !== 'undefined') {
  const savedInstanceID = readSavedInstanceID()
  if (savedInstanceID) {
    connect(savedInstanceID)
  }
}

export { getSocketUrl, connect, disconnect, subscribe }
