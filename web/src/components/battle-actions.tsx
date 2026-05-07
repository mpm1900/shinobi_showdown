import { ActionControl } from './action-control'
import { useStore } from '@tanstack/react-store'
import { gameStore } from '#/lib/stores/game'
import type { Actor } from '#/lib/game/actor'
import { battleContext, setContext } from '#/lib/stores/battle-context'
import { AnimatePresence, LayoutGroup, motion } from 'motion/react'
import { BattleCards } from './battle-cards'

import { memo, useMemo } from 'react'
import { contextToString } from '#/lib/game/context'

const BattleActions = memo(function BattleActions({ actor }: { actor: Actor }) {
  const status = useStore(gameStore, (g) => g.status)
  const actions = useStore(gameStore, (g) => g.actions)
  const context_key = useStore(battleContext, (c) => contextToString(c))
  const context = useMemo(() => {
    const c = battleContext.state
    return {
      action_ID: c.action_ID,
      source_player_ID: c.source_player_ID,
      parent_actor_ID: c.parent_actor_ID,
      source_actor_ID: c.source_actor_ID,
      target_actor_IDs: c.target_actor_IDs,
      target_position_IDs: c.target_position_IDs,
    }
  }, [context_key])
  const action = actor.actions.find((a) => a.ID === context.action_ID)
  const idle = status === 'idle'
  const staged = actions.find((tx) => tx.context.source_actor_ID === actor.ID)

  return (
    <LayoutGroup id="battle-actions">
      <div className="pointer-events-none relative flex w-full flex-col items-center gap-4 pb-8">
        <div className="pointer-events-auto">
          <AnimatePresence mode="wait" initial={false}>
            {action && idle && (
              <motion.div
                key={`${actor.ID}.${action.ID}`}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                exit={{ opacity: 0 }}
                transition={{
                  type: 'spring',
                  stiffness: 520,
                  damping: 36,
                  mass: 0.6,
                }}
              >
                {!staged && (
                  <div className="grid place-items-center mb-6 nanum-brush-script-regular text-5xl text-shadow-[2px_2px_0px_#000000]">
                    <div>
                      {action ? action.config.name : 'Choose an Action'}
                    </div>
                  </div>
                )}
                <ActionControl
                  action={action}
                  staged={staged}
                  enabled={idle && !!actor.position_ID && !action?.disabled}
                  context={context}
                  onContextChange={setContext}
                />
              </motion.div>
            )}
          </AnimatePresence>
        </div>
        <BattleCards actor={actor} />
      </div>
    </LayoutGroup>
  )
})

export { BattleActions }
