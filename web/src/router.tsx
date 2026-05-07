import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { setupRouterSsrQueryIntegration } from '@tanstack/react-router-ssr-query'
import { routeTree } from './routeTree.gen'

import { getContext } from './integrations/tanstack-query/root-provider'

export function getRouter() {
  const queryContext = getContext()
  const router = createTanStackRouter({
    routeTree,

    context: {
      ...queryContext,
      auth: {
        user: null,
      },
    },

    scrollRestoration: true,
    defaultPreload: 'intent',
    defaultPreloadStaleTime: 0,
  })

  setupRouterSsrQueryIntegration({
    router,
    queryClient: queryContext.queryClient,
  })

  return router
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
