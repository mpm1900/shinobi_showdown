import { getApiBaseUrl } from '#/lib/server/api-base'
import { queryOptions } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import type { Action } from '../game/action'

const getActions = createServerFn().handler(async () => {
  const base = getApiBaseUrl()
  const response = await fetch(`${base}/api/actions`)
  const data = await response.json()
  return data as Array<Action>
})

const actionsQuery = queryOptions({
  queryKey: ['actions'],
  queryFn: () => getActions(),
  staleTime: 60000,
  gcTime: 60000,
})

export { actionsQuery, getActions }
