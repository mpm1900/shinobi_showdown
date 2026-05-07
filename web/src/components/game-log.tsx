import { useEffect, useRef } from 'react'
import { type GameLog as GameLogType } from '#/lib/game/log'
import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { clientsStore } from '#/lib/stores/clients'
import { cn } from '#/lib/utils'
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from './ui/accordion'
import { ScrollArea } from './ui/scroll-area'

function GameLogItem({
  item,
  clientID,
}: {
  item: GameLogType
  clientID: string
}) {
  const actors = useStore(gameStore, (g) => g.actors)
  const source = actors.find((a) => a.ID === item.context.source_actor_ID)
  const action = source?.actions?.find((a) => a.ID === item.context.action_ID)

  const sourceText = source?.name ?? 'Unknown'
  const actionText = action?.config.name ?? 'Unknown action'

  const tokens = item.text.split(/(\$source\$|\$action\$)/g)

  return (
    <span>
      {tokens.map((token, idx) => {
        if (!token) return null
        if (token === '$source$')
          return (
            <span
              key={`source-${idx}`}
              className={cn({
                'text-blue-300': clientID === source?.player_ID,
                'text-red-300': clientID !== source?.player_ID,
              })}
            >
              {sourceText}
            </span>
          )
        if (token === '$action$')
          return (
            <span key={`action-${idx}`} className="text-foreground">
              {actionText}
            </span>
          )
        return <span key={`text-${idx}`}>{token}</span>
      })}
    </span>
  )
}

function GameLogList() {
  const log = useStore(gameStore, (s) => s.log)
  const clientID = useStore(clientsStore, (s) => s.me?.ID ?? '')
  const endRef = useRef<HTMLLIElement | null>(null)

  const lastLogID = log.length > 0 ? log[log.length - 1].ID : null

  useEffect(() => {
    endRef.current?.scrollIntoView({ block: 'end' })
  }, [lastLogID])

  return (
    <ul>
      {log.map((item, index) => (
        <li
          key={item.ID}
          ref={index === log.length - 1 ? endRef : undefined}
          className="text-muted-foreground"
        >
          <GameLogItem item={item} clientID={clientID} />
        </li>
      ))}
    </ul>
  )
}

function GameLog() {
  return (
    <Accordion
      defaultValue={['log']}
      type="multiple"
      className="bg-stone-950/80 px-3 rounded-sm border border-stone-300/30 ring-1 ring-black min-w-96 mt-4"
    >
      <AccordionItem value="log">
        <AccordionTrigger>Log</AccordionTrigger>
        <AccordionContent>
          <ScrollArea className="h-40">
            <GameLogList />
          </ScrollArea>
        </AccordionContent>
      </AccordionItem>
    </Accordion>
  )
}

export { GameLogItem, GameLogList, GameLog }
