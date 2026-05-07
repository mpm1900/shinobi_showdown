import { battleContext } from "#/lib/stores/battle-context"
import { gameStore } from "#/lib/stores/game"
import { useStore } from "@tanstack/react-store"

function useActiveActor() {
  const source_actor_ID = useStore(battleContext, (c) => c.source_actor_ID)
  const actors = useStore(gameStore, g => g.actors)
  const actor = actors.find((a) => a.ID === source_actor_ID)
  return actor
}

export { useActiveActor }
