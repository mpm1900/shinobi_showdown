import { AnimatePresence, motion, type Transition } from 'motion/react'
import { cn } from '#/lib/utils'
import { battleContext, setContextAction } from '#/lib/stores/battle-context'
import { useStore } from '@tanstack/react-store'
import { gameStore } from '#/lib/stores/game'
import type { Actor } from '#/lib/game/actor'
import { ActionCard } from './action-card'

const collection_transition: Transition = {
  type: 'spring',
  stiffness: 800,
  damping: 50,
  mass: 0.68,
}
const card_transition: Transition = {
  type: 'spring',
  stiffness: 760,
  damping: 38,
  mass: 0.5,
}

function BattleCards({ actor }: { actor: Actor }) {
  const queued_actions = useStore(gameStore, (g) => g.queued_actions)
  const prompt = useStore(gameStore, g => g.prompt)
  const status = useStore(gameStore, (g) => g.status)
  const actions = useStore(gameStore, (g) => g.actions)
  const action_ID = useStore(battleContext, (c) => c.action_ID)
  const player = useStore(gameStore, (g) =>
    g.players.find((p) => p.ID === actor.player_ID)
  )
  const has_queued_action = queued_actions[actor.ID]
  const idle = status === 'idle' && !prompt
  const staged_action = actions.find(
    (tx) => tx.context.source_actor_ID === actor.ID
  )
  return (
    <div className="pointer-events-none fixed inset-x-0 bottom-12 z-0 flex justify-center">
      <AnimatePresence mode="wait" initial={false}>
        <motion.div
          key={actor.ID}
          initial={{ y: 100, opacity: 0 }}
          animate={{ y: 0, opacity: idle && !has_queued_action ? 1 : 0 }}
          exit={{ y: 100, opacity: 0 }}
          transition={collection_transition}
          className="flex items-end -space-x-11 px-4 transform-gpu"
        >
          {actor.actions.map((a, i) => {
            const selected = action_ID === a.ID
            const center = (actor.actions.length - 1) / 2
            const fanRotate = (i - center) * 1.25

            return (
              <motion.div
                key={`${actor.ID}:${a.ID}`}
                className={cn(
                  'relative pointer-events-auto transform-gpu will-change-transform',
                  {
                    'pointer-events-none': !!staged_action && !selected,
                  }
                )}
                style={{
                  zIndex: selected ? 80 : 20 + i,
                  willChange: 'transform',
                }}
                initial={{ y: 12, scale: 0.985 }}
                animate={{
                  y: selected ? -64 : 0,
                  scale: selected ? 1.015 : 1,
                  rotate: selected ? 0 : fanRotate,
                }}
                exit={{ y: 10, scale: 0.985 }}
                whileHover={{
                  y: selected ? -64 : -36,
                  scale: selected ? 1.05 : 1.02,
                  rotate: 0,
                  zIndex: 90,
                  paddingBottom: 28,
                  marginBottom: -28,
                }}
                whileTap={{ scale: 1.01 }}
                transition={card_transition}
              >
                <ActionCard
                  action={a}
                  usedSummon={!!player?.used_summon}
                  disabled={(!!staged_action && !selected) || a.disabled}
                  selected={selected}
                  onClick={() => setContextAction(a.ID)}
                  style={{
                    zIndex: selected ? 80 : 20 + i,
                  }}
                />
              </motion.div>
            )
          })}
        </motion.div>
      </AnimatePresence>
    </div>
  )
}

export { BattleCards }
