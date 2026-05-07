import type { Context } from '#/lib/game/context'
import { subscribe } from '#/lib/socket/connect'
import { clientsStore } from '#/lib/stores/clients'
import { sendContextMessage } from '#/lib/stores/socket'
import { useStore } from '@tanstack/react-store'
import { useEffect, useRef, useState } from 'react'

type GetTargetsState = {
  context: Context | null
  loading: boolean
}

function contextScopeKey(context: Context): string {
  const targetPositionIDs = context.target_position_IDs ?? []
  return [
    context.action_ID ?? '',
    context.source_player_ID ?? '',
    context.parent_actor_ID ?? '',
    context.source_actor_ID ?? '',
    targetPositionIDs.join('|'),
  ].join('::')
}

function useGetTargets(context: Context, prompt_ID?: string) {
  const client = useStore(clientsStore, (c) => c.me)

  const [state, setState] = useState<GetTargetsState>({
    context: null,
    loading: false,
  })

  const latestScopeRef = useRef<string>('')
  const latestRequestRef = useRef(0)
  const latestContextRef = useRef(context)

  const requestScope = contextScopeKey(context)

  useEffect(() => {
    latestContextRef.current = context
  }, [context])

  useEffect(() => {
    return subscribe((_, message) => {
      if (message?.type !== 'target-IDs' || !message.context) return

      // Ignore responses from unrelated get-targets requests.
      const responseScope = contextScopeKey(message.context)
      if (responseScope !== latestScopeRef.current) return

      setState({
        context: message.context,
        loading: false,
      })
    })
  }, [])

  useEffect(() => {
    if (!client) return

    const requestID = latestRequestRef.current + 1
    latestRequestRef.current = requestID
    latestScopeRef.current = requestScope

    setState((s) => ({
      ...s,
      loading: true,
    }))

    const sent = sendContextMessage({
      type: 'get-targets',
      context: latestContextRef.current,
      client_ID: client.ID,
      prompt_ID,
    })

    if (!sent && latestRequestRef.current === requestID) {
      setState((s) => ({
        ...s,
        loading: false,
      }))
    }
  }, [client?.ID, prompt_ID, requestScope])

  return state
}

export { useGetTargets }
