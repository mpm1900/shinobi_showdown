import { contextToString, type Context } from '#/lib/game/context'
import { subscribe } from '#/lib/socket/connect'
import { clientsStore } from '#/lib/stores/clients'
import { sendContextMessage } from '#/lib/stores/socket'
import { useStore } from '@tanstack/react-store'
import { useEffect, useRef, useState } from 'react'

type ValidateContextState = {
  valid: boolean
  loading: boolean
}

type PendingRequest = {
  seq: number
  contextKey: string
}

function useValidateContext(context: Context, prompt_ID?: string) {
  const client = useStore(clientsStore, (c) => c.me)
  const [state, setState] = useState<ValidateContextState>({
    valid: false,
    loading: false,
  })

  const pendingRef = useRef<PendingRequest | null>(null)
  const seqRef = useRef(0)

  useEffect(() => {
    return subscribe((_, message) => {
      if (message?.type !== 'validate-context' || !message.context) return

      const pending = pendingRef.current
      if (!pending) return

      const responseContextKey = contextToString(message.context)
      if (responseContextKey !== pending.contextKey) return

      setState({
        valid: message.valid ?? false,
        loading: false,
      })

      pendingRef.current = null
    })
  }, [])

  const contextKey = contextToString(context)

  useEffect(() => {
    if (!client) return

    const seq = ++seqRef.current

    pendingRef.current = { seq, contextKey }
    setState((s) => ({ ...s, loading: true }))

    sendContextMessage({
      type: 'validate-context',
      context,
      client_ID: client.ID,
      prompt_ID,
    })
  }, [client?.ID, contextKey, prompt_ID])

  return state
}

export { useValidateContext }
