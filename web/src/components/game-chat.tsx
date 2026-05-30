import { ArrowRight } from 'lucide-react'
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from './ui/input-group'
import { sendContextMessage } from '#/lib/stores/socket'
import { useStore } from '@tanstack/react-store'
import { clientsStore } from '#/lib/stores/clients'
import { NULL_CONTEXT } from '#/lib/game/context'
import { gameStore } from '#/lib/stores/game'
import { useForm } from '@tanstack/react-form'
import z from 'zod'
import { chatsStore } from '#/lib/stores/chat'
import { ScrollArea } from './ui/scroll-area'
import { useEffect, useRef } from 'react'

const ChatSchema = z.object({
  text: z.string().min(1),
})

function GameChat() {
  const client = useStore(clientsStore, (state) => state.me!)
  const player = useStore(gameStore, (g) =>
    g.players.find((p) => p.ID === client.ID)
  )
  const chats = useStore(chatsStore, (c) => c)
  const form = useForm({
    defaultValues: {
      text: '',
    },
    validators: {
      onChange: ChatSchema,
    },
    onSubmit: async ({ value, formApi }) => {
      sendContextMessage({
        type: 'send-chat',
        client_ID: client.ID,
        context: NULL_CONTEXT,
        chat_message: {
          client_ID: client.ID,
          text: value.text,
          timestamp: new Date().toISOString(),
          username: player?.user.email ?? client.ID,
        },
      })
      formApi.reset()
    },
  })

  const endRef = useRef<HTMLLIElement | null>(null)
  const lastChatID = chats.length > 0 ? chats[chats.length - 1].ID : null

  useEffect(() => {
    endRef.current?.scrollIntoView({ block: 'end' })
  }, [lastChatID])

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault()
        e.stopPropagation()
        form.handleSubmit()
      }}
    >
      <div className="flex flex-col h-40 gap-2">
        <ScrollArea className="h-28">
          <ul className="flex-1 flex flex-col">
            {chats.map((chat, index) => (
              <li
                key={chat.ID}
                ref={index === chats.length - 1 ? endRef : undefined}
                className="space-x-2 inline-block"
              >
                <span className="font-bold text-amber-500">
                  {chat.username}
                </span>
                <span>{chat.text}</span>
              </li>
            ))}
          </ul>
        </ScrollArea>
        <div>
          <form.Field name="text">
            {(field) => (
              <InputGroup>
                <InputGroupInput
                  placeholder="Send a message"
                  value={field.state.value}
                  onChange={(e) => field.handleChange(e.target.value)}
                />
                <InputGroupAddon align="inline-end">
                  <InputGroupButton
                    variant="default"
                    disabled={!field.form.state.canSubmit || field.form.state.isPristine}
                    type="submit"
                  >
                    <ArrowRight />
                  </InputGroupButton>
                </InputGroupAddon>
              </InputGroup>
            )}
          </form.Field>
        </div>
      </div>
    </form>
  )
}

export { GameChat }
