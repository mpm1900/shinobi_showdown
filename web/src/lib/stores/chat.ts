import { Store } from '@tanstack/store'
import type { ChatMessage } from '../socket/request'

const chatsStore = new Store<ChatMessage[]>([])

function pushChatMessage(message: ChatMessage) {
  chatsStore.setState((prev) => [...prev, message])
}

export { chatsStore, pushChatMessage }
