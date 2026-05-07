import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { ActorCard } from './actor-card'
import { clientsStore } from '#/lib/stores/clients'
import { Item } from './ui/item'
import { cn } from '#/lib/utils'
import { motion, AnimatePresence } from 'motion/react'
import { battleContext } from '#/lib/stores/battle-context'

function PlayerPositions({
  flip,
  player_ID,
  selected,
  onSelectedChange,
}: {
  flip: boolean
  player_ID: string
  selected?: string
  onSelectedChange?: (selected: string) => void
}) {
  const game = useStore(gameStore, (g) => g)
  const client = useStore(clientsStore, (c) => c.me)
  const bc = useStore(battleContext, (c) => c)
  const player = game.players.find((p) => p.ID === player_ID)
  const coef = flip ? -1 : 1
  const is_running = game.status === 'running'
  const running_context = is_running ? game.active_transaction?.context : null
  const targetedActorIDs = new Set<string>([
    ...(bc.hover_target_IDs ?? []),
    ...(bc.target_actor_IDs ?? []),
    ...(running_context?.target_actor_IDs ?? []),
  ])
  const targetedPositionIDs = new Set<string>([
    ...(bc.target_position_IDs ?? []),
    ...(running_context?.target_position_IDs ?? []),
  ])

  if (!player) return null

  return (
    <div className="flex gap-8">
      {player.positions.map((pos) => {
        const actor_ID = pos.actor_ID
        const targeted =
          (!!actor_ID && targetedActorIDs.has(actor_ID)) ||
          targetedPositionIDs.has(pos.ID)
        const is_source = running_context?.source_actor_ID === actor_ID
        const is_selected = selected === actor_ID
        return (
          <div
            key={pos.ID}
            className={cn('flex items-end', flip && 'items-start')}
          >
            <AnimatePresence mode="wait" initial={false}>
              <motion.div
                key={pos.actor_ID ?? `empty:${pos.ID}`}
                className="transform-gpu"
                initial={{ y: 24 * coef, opacity: 0 }}
                animate={{
                  y: 0,
                  opacity: 1,
                  scale: is_selected || is_source ? 1.15 : 1,
                }}
                exit={{ y: 24 * coef, opacity: 0 }}
                transition={{
                  type: 'spring',
                  stiffness: 520,
                  damping: 36,
                  mass: 0.6,
                }}
              >
                {pos.actor_ID ? (
                  <ActorCard
                    actor={game.actors.find((a) => a.ID === pos.actor_ID)}
                    client_ID={client?.ID}
                    selected={is_selected || is_source}
                    source={is_source}
                    targeted={targeted}
                    data-actor-anchor={pos.actor_ID}
                    onClick={() => onSelectedChange?.(pos.actor_ID!)}
                    className={flip ? 'flex-col-reverse' : ''}
                    summonClass={!flip ? 'top-auto! bottom-0!' : ''}
                  />
                ) : (
                  <Item className="py-6 w-108">
                    <span className="text-center w-full text-muted-foreground">
                      Empty
                    </span>
                  </Item>
                )}
              </motion.div>
            </AnimatePresence>
          </div>
        )
      })}
    </div>
  )
}

export { PlayerPositions }
