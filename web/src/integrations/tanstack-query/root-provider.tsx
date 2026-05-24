import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import type { ReactNode } from 'react'

function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
      },
    },
  })
}

let context:
  | {
      queryClient: QueryClient
    }
  | undefined

export function getContext() {
  // Never share query cache between server requests.
  if (typeof window === 'undefined') {
    return {
      queryClient: createQueryClient(),
    }
  }

  if (context) {
    return context
  }

  context = {
    queryClient: createQueryClient(),
  }

  return context
}

export default function TanStackQueryProvider({
  children,
}: {
  children: ReactNode
}) {
  const { queryClient } = getContext()

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  )
}
