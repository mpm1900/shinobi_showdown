import {
  mutationOptions,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import z from 'zod'
import type { User } from '#/lib/queries/auth'
import { setResponseCookie } from '#/utils/set-cookie'
import { getApiBaseUrl } from '#/lib/server/api-base'

const requestSchema = z.object({
  email: z.string(),
  password: z.string(),
  secret: z.string(),
})

const signup = createServerFn({ method: 'POST' })
  .inputValidator(requestSchema)
  .handler(async ({ data }) => {
    const response = await fetch(`${getApiBaseUrl()}/api/auth/signup`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })

    setResponseCookie(response)

    if (!response.ok) {
      throw new Error(`Signup failed with status ${response.status}`)
    }

    return (await response.json()) as User
  })

function useSignup() {
  const queryClient = useQueryClient()
  return useMutation(
    mutationOptions({
      mutationKey: ['signup'],
      mutationFn: async (data: z.output<typeof requestSchema>) => {
        return await signup({ data })
      },
      onSuccess: (user) => {
        queryClient.setQueryData(['me'], user)
        queryClient.removeQueries({ queryKey: ['teams'] })
      },
    })
  )
}

export { useSignup }
