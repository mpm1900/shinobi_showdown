import { getApiBaseUrl } from '#/lib/server/api-base'
import { mutationOptions, useMutation } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import { getRequest } from '@tanstack/react-start/server'
import z from 'zod'
import type { Team } from '../queries/teams'
import { TeamSchema } from '../stores/config'

const upsertTeamSchema = z.object({
  id: z.string().nullable(),
  config: TeamSchema,
})

const upsertTeam = createServerFn({ method: 'POST' })
  .inputValidator(upsertTeamSchema)
  .handler(async ({ data }) => {
    const request = getRequest()
    const cookies = request?.headers.get('cookie') || ''

    const response = await fetch(`${getApiBaseUrl()}/api/teams/upsert`, {
      method: 'POST',
      body: JSON.stringify(data),
      headers: {
        'Content-Type': 'application/json',
        Cookie: cookies,
      },
    })

    const team = await response.json()
    return team as Team
  })

function useUpsertTeam() {
  return useMutation(
    mutationOptions({
      mutationKey: ['upsert-team'],
      mutationFn: async (data: z.output<typeof upsertTeamSchema>) => {
        return await upsertTeam({ data })
      },
    })
  )
}

export { useUpsertTeam }
