import { queryOptions, useQuery } from '@tanstack/react-query'
import { createServerFn } from '@tanstack/react-start'
import { getRequest } from '@tanstack/react-start/server'

export type User = {
  id: string
  username: string
  email: string
  created_at: string
}

const getMe = createServerFn().handler(async () => {
  const request = getRequest()
  const cookies = request?.headers.get('cookie') || ''

  const response = await fetch(`${process.env.API_URL}/api/auth/me`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      Cookie: cookies,
    },
  })

  if (!response.ok) {
    return null
  }

  return (await response.json()) as User
})

export const meQuery = queryOptions({
  queryKey: ['me'],
  queryFn: () => getMe(),
})

export function useUser() {
  return useQuery(meQuery)
}
