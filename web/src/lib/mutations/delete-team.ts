import { createServerFn } from '@tanstack/react-start'
import z from 'zod'
import { mutationOptions, useMutation } from '@tanstack/react-query'
import { getRequest } from '@tanstack/react-start/server'

const deleteTeamSchema = z.object({
  team_id: z.string(),
})

const deleteTeam = createServerFn({ method: 'POST' })
  .inputValidator(deleteTeamSchema)
  .handler(async ({ data }) => {
    const request = getRequest()
    const cookies = request?.headers.get('cookie') || ''

    await fetch(`${process.env.API_URL}/api/teams/${data.team_id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        Cookie: cookies,
      },
    })
  })

function useDeleteTeam() {
  return useMutation(
    mutationOptions({
      mutationKey: ['delete-team'],
      mutationFn: async (team_id: string | null) => {
        if (team_id === null) return
        await deleteTeam({
          data: {
            team_id,
          },
        })
      },
    })
  )
}

export { useDeleteTeam }
