import {
  mutationOptions,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import z from 'zod'
import type { User } from '../queries/auth'
import { setResponseCookie } from '#/utils/set-cookie'
import { getApiBaseUrl } from '#/lib/server/api-base'

const requestSchema = z.object({
  email: z.string(),
  password: z.string(),
})

const login = createServerFn({ method: 'POST' })
  .inputValidator(requestSchema)
  .handler(async ({ data }) => {
    const response = await fetch(`${getApiBaseUrl()}/api/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      throw new Error(`Login failed with status ${response.status}`)
    }

    setResponseCookie(response)

    return (await response.json()) as User
  })

function useLogin() {
  const queryClient = useQueryClient()
  return useMutation(
    mutationOptions({
      mutationKey: ['login'],
      mutationFn: async (data: z.output<typeof requestSchema>) => {
        return await login({ data })
      },
      onSuccess: (user) => {
        queryClient.setQueryData(['me'], user)
        queryClient.removeQueries({ queryKey: ['teams'] })
      },
    })
  )
}

export { useLogin }
