import { setContextPlayer } from '../stores/battle-context'
import { pushChatMessage } from '../stores/chat'
import { clientsStore } from '../stores/clients'
import { gameStore } from '../stores/game'
import type { SocketResponse } from './request'

function messageReducer(message: SocketResponse | null) {
  if (!message?.type) return
  switch (message.type) {
    case 'game': {
      if (message.state) {
        gameStore.setState(() => message.state!)
      }
      return
    }
    case 'clients': {
      clientsStore.setState((c) => ({
        ...c,
        clients: message.clients!,
      }))
      return
    }
    case 'join-success': {
      if (message.state) {
        gameStore.setState(() => message.state!)
      }
      clientsStore.setState((c) => ({
        ...c,
        me: message.clients![0],
      }))
      setContextPlayer(message.clients![0].ID)
      return
    }
    case 'new-chat': {
      if (message.chat_message) {
        pushChatMessage(message.chat_message)
      }
    }
  }
}

export { messageReducer }
