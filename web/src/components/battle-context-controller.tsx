import type { ActionTransaction } from '#/lib/game/action'
import { battleContext, setContextSource } from '#/lib/stores/battle-context'
import { gameStore } from '#/lib/stores/game'
import { useStore } from '@tanstack/react-store'
import { useEffect } from 'react'

function hasAction(a_id: string, actions: ActionTransaction[]) {
  return actions.find((tx) => tx.context.source_actor_ID === a_id)
}

function BattleContextController() {
  const game = useStore(gameStore, (g) => g)
  const source_player_ID = useStore(battleContext, (c) => c.source_player_ID)
  const actors = game.actors.filter(
    (a) => a.player_ID === source_player_ID
  )
  const active_actors = actors.filter((a) => !!a.position_ID)
  const actionable_actors = active_actors.filter(
    (a) => !a.stunned && !hasAction(a.ID, game.actions)
  )

  useEffect(() => {
    if (game.turn.phase === 'main') {
      const first = actionable_actors[0]
      if (!first) return

      setContextSource(first.ID, game)
    }
  }, [game.turn.phase, actionable_actors.length])

  return null
}

export { BattleContextController }
