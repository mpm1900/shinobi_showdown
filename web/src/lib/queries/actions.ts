import { createServerFn } from '@tanstack/react-start'
import { queryOptions } from '@tanstack/react-query'
import type { Action } from '../game/action'

const getActions = createServerFn().handler(async () => {
  const response = await fetch(`${process.env.API_URL}/api/actions`)
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
