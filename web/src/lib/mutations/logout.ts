import { getApiBaseUrl } from '#/lib/server/api-base'
import { setResponseCookie } from '#/utils/set-cookie'
import {
  mutationOptions,
  useMutation,
  useQueryClient,
} from '@tanstack/react-query'
import { useRouter } from '@tanstack/react-router'
import { createServerFn } from '@tanstack/react-start'
import { getRequest } from '@tanstack/react-start/server'

const logout = createServerFn({ method: 'POST' }).handler(async () => {
  const request = getRequest()
  const cookies = request?.headers.get('cookie') || ''

  const response = await fetch(`${getApiBaseUrl()}/api/auth/logout`, {
    method: 'POST',
    headers: {
      Cookie: cookies,
    },
  })

  setResponseCookie(response)
})

function useLogout() {
  const queryClient = useQueryClient()
  const router = useRouter()

  return useMutation(
    mutationOptions({
      mutationKey: ['logout'],
      mutationFn: async () => {
        await logout()
      },
      onSuccess: () => {
        queryClient.setQueryData(['me'], null)
        queryClient.removeQueries({ queryKey: ['teams'] })
        router.navigate({ to: '/login' })
      },
    })
  )
}

export { useLogout }
