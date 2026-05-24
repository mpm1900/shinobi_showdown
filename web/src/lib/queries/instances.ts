import { getApiBaseUrl } from '#/lib/server/api-base'
import { queryOptions } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'

type Instance = {
  ID: string
}

const getInstances = createServerFn().handler(async () => {
  const response = await fetch(`${getApiBaseUrl()}/api/instances`)
  const data = await response.json()
  return data as Instance[]
})

const instancesQuery = queryOptions({
  queryKey: ['instances'],
  queryFn: () => getInstances(),
})

export { getInstances, instancesQuery }
