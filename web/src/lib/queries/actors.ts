import { createServerFn } from '@tanstack/react-start'
import type { ActorDef } from '../game/actor'
import { queryOptions } from '@tanstack/react-query'

const getActors = createServerFn().handler(async () => {
  const response = await fetch(`${process.env.API_URL}/api/actors`)
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
