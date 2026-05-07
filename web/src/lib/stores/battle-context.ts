import { Store } from '@tanstack/store'
import { NULL_CONTEXT, type Context } from '../game/context'
import type { Game } from '../game/game'
import type { Actor } from '../game/actor'

type BattleContextState = Context & {
  previous_action_IDs: Record<string, string>
  hover_target_IDs: Array<string>
}

const battleContext = new Store<BattleContextState>({
  ...NULL_CONTEXT,
  previous_action_IDs: {},
  hover_target_IDs: []
})

function setContextPlayer(player_ID: string) {
  battleContext.setState((s) => ({
    ...s,
    source_player_ID: player_ID,
    parent_actor_ID: null,
    source_actor_ID: null,
    target_actor_IDs: [],
    target_position_IDs: [],
    previous_action_IDs: {},
    hover_target_IDs: [],
  }))
}

function clearHoverTargets() {
  battleContext.setState((s) => ({
    ...s,
    hover_target_IDs: [],
  }))
}
function addHoverTarget(ID: string) {
  if (!ID) return

  battleContext.setState((s) => ({
    ...s,
    hover_target_IDs: s.hover_target_IDs.includes(ID)
      ? s.hover_target_IDs
      : s.hover_target_IDs.concat(ID),
  }))
}
function removeHoverTarget(ID: string) {
  battleContext.setState((s) => ({
    ...s,
    hover_target_IDs: s.hover_target_IDs.filter((t) => t !== ID),
  }))
}

function getNextActionableActor(
  game: Game,
  context: BattleContextState,
  actor_ID: string
) {
  const playerID = context.source_player_ID
  if (!playerID) return null

  const actedActorIDs = new Set(
    game.actions
      .map((tx) => tx.context.source_actor_ID)
      .filter(Boolean) as string[]
  )
  actedActorIDs.add(actor_ID)

  const actors = game.actors.filter((a) => a.player_ID === playerID)
  const activeActors = actors.filter((a) => !!a.position_ID)
  const actionableActors = activeActors.filter(
    (a) => !a.stunned && !actedActorIDs.has(a.ID)
  )

  return actionableActors[0] ?? null
}

function nextAction(actor: Actor, previous_action_IDs: {
  [x: string]: string;
}) {

  const preferredActionID = (previous_action_IDs[actor.ID] as string | null) ?? null
  const preferredAction = actor.actions.find((a) => a.ID === preferredActionID)
  const fallbackActionID = actor.actions.find(a => !a.disabled)?.ID ?? null
  let nextActionID = preferredActionID ?? fallbackActionID
  if (preferredAction?.disabled === true) {
    nextActionID = fallbackActionID
  }

  return nextActionID
}

function setActionID(actor_ID: string, action_ID: string, game: Game) {
  battleContext.setState((s) => {
    const previous_action_IDs = {
      ...s.previous_action_IDs,
      [actor_ID]: action_ID,
    }

    let resolvedNextActor = getNextActionableActor(game, s, actor_ID)
    let nextActionID = resolvedNextActor
      ? nextAction(resolvedNextActor, previous_action_IDs)
      : null

    if (!resolvedNextActor || !nextActionID) {
      const playerID =
        s.source_player_ID ??
        game.actors.find((a) => a.ID === actor_ID)?.player_ID ??
        null

      if (playerID) {
        const actedActorIDs = new Set(
          game.actions
            .map((tx) => tx.context.source_actor_ID)
            .filter(Boolean) as string[]
        )
        actedActorIDs.add(actor_ID)

        const fallbackActors = game.actors.filter(
          (a) =>
            a.player_ID === playerID &&
            !!a.position_ID &&
            !a.stunned &&
            !actedActorIDs.has(a.ID)
        )

        const orderedFallbackActors = resolvedNextActor
          ? [
              resolvedNextActor,
              ...fallbackActors.filter((a) => a.ID !== resolvedNextActor?.ID),
            ]
          : fallbackActors

        for (const candidate of orderedFallbackActors) {
          const candidateActionID = nextAction(candidate, previous_action_IDs)
          if (candidateActionID) {
            resolvedNextActor = candidate
            nextActionID = candidateActionID
            break
          }
        }
      }
    }

    if (!resolvedNextActor || !nextActionID) {
      return {
        ...s,
        action_ID: null,
        parent_actor_ID: null,
        source_actor_ID: null,
        target_actor_IDs: [],
        target_position_IDs: [],
        previous_action_IDs,
        hover_target_IDs: [],
      }
    }

    return {
      ...s,
      action_ID: nextActionID,
      parent_actor_ID: resolvedNextActor.ID,
      source_actor_ID: resolvedNextActor.ID,
      target_actor_IDs: [],
      target_position_IDs: [],
      previous_action_IDs,
      hover_target_IDs: [],
    }
  })
}

function setContextSource(source_ID: string, game: Game) {
  const existing = game.actions.find(
    (t) => t.context.source_actor_ID === source_ID
  )

  if (existing) {
    return battleContext.setState((s) => ({
      ...s,
      ...existing.context,
      hover_target_IDs: [],
    }))
  }

  const source = game.actors.find((a) => a.ID === source_ID)
  if (!source) return

  const prev_ID = nextAction(source, battleContext.state.previous_action_IDs)
  battleContext.setState((s) => ({
    ...s,
    action_ID: prev_ID,
    parent_actor_ID: source_ID,
    source_actor_ID: source_ID,
    target_actor_IDs: [],
    target_position_IDs: [],
    hover_target_IDs: [],
  }))
}

function setContextAction(action_ID: string) {
  battleContext.setState((s) => ({
    ...s,
    action_ID,
    target_actor_IDs: [],
    target_position_IDs: [],
    hover_target_IDs: [],
  }))
}

function setContext(context: Context) {
  battleContext.setState((s) => ({
    ...s,
    ...context,
    hover_target_IDs: [],
  }))
}

export {
  battleContext,
  setContext,
  setContextPlayer,
  setActionID,
  setContextSource,
  setContextAction,
  clearHoverTargets,
  addHoverTarget,
  removeHoverTarget,
}
