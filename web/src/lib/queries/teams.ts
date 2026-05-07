import { createServerFn } from '@tanstack/react-start'
import { getRequest } from '@tanstack/react-start/server'
import { queryOptions } from '@tanstack/react-query'
import type { TeamConfig } from '../stores/config'

type Team = {
  id: string | null
  team_config: TeamConfig
  created_at: string | null
}

const getTeams = createServerFn().handler(async () => {
  const request = getRequest()
  const cookies = request?.headers.get('cookie') || ''

  const response = await fetch(`${process.env.API_URL}/api/teams`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Cookie: cookies,
    },
  })

  if (!response.ok) {
    return []
  }

  const data = await response.json()
  return data as Team[]
})

const teamsQuery = queryOptions({
  queryKey: ['teams'],
  queryFn: () => getTeams(),
  staleTime: 60000,
  gcTime: 60000,
})

export type { Team }
export { teamsQuery, getTeams }
