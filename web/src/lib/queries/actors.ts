import { getApiBaseUrl } from '#/lib/server/api-base'
import { queryOptions } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import type { ActorDef } from '../game/actor'

const getActors = createServerFn().handler(async () => {
  const response = await fetch(`${getApiBaseUrl()}/api/actors`)
  const data = await response.json()
  return data as Array<ActorDef>
})

const actorsQuery = queryOptions({
  queryKey: ['actors'],
  queryFn: () => getActors(),
  staleTime: 60000,
  gcTime: 60000,
})

export { actorsQuery, getActors }
