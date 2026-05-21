import {
  ClientOnly,
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from '@tanstack/react-router'
import TanStackQueryProvider from '../integrations/tanstack-query/root-provider'

import styles from '../styles.css?url'

import type { QueryClient } from '@tanstack/react-query'
import { TooltipProvider } from '#/components/ui/tooltip'
import type { User } from '#/lib/queries/auth'
import { VantaBackground } from '#/components/vanta-background'

interface RouterContext {
  queryClient: QueryClient
  auth: {
    user: User | null
  }
}

// const THEME_INIT_SCRIPT = `(function(){try{var stored=window.localStorage.getItem('theme');var mode=(stored==='light'||stored==='dark'||stored==='auto')?stored:'auto';var prefersDark=window.matchMedia('(prefers-color-scheme: dark)').matches;var resolved=mode==='auto'?(prefersDark?'dark':'light'):mode;var root=document.documentElement;root.classList.remove('light','dark');root.classList.add(resolved);if(mode==='auto'){root.removeAttribute('data-theme')}else{root.setAttribute('data-theme',mode)}root.style.colorScheme=resolved;}catch(e){}})();`

import { meQuery } from '#/lib/queries/auth'
import { useStore } from '@tanstack/react-store'
import { uiStore } from '#/lib/stores/ui'

export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async ({ context, location }) => {
    if (location.pathname === '/up') {
      return { auth: { user: null } }
    }
    const user = await context.queryClient.fetchQuery(meQuery)
    return {
      auth: {
        user,
      },
    }
  },
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Shinobi Showdown',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: styles,
      },
      {
        rel: 'preconnect',
        href: 'https://fonts.googleapis.com',
      },
      {
        rel: 'preconnect',
        href: 'https://fonts.gstatic.com',
        crossOrigin: 'anonymous',
      },
      {
        rel: 'stylesheet',
        href: 'https://fonts.googleapis.com/css2?family=Nanum+Brush+Script&family=Yeon+Sung&display=swap',
      },
    ],
    scripts: [{ src: '/scripts/three.js' }, { src: '/scripts/fog.js' }],
  }),
  shellComponent: RootDocument,
  errorComponent: RootErrorComponent,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  const bgEnabled = useStore(uiStore, (s) => s.bgEnabled)

  return (
    <html lang="en" className="dark" suppressHydrationWarning>
      <head>
        {/*<script dangerouslySetInnerHTML={{ __html: THEME_INIT_SCRIPT }} />*/}

        <HeadContent />
      </head>
      <body className="font-sans antialiased wrap-anywhere overflow-x-hidden flex flex-col bg-stone-800">
        {bgEnabled && (
          <ClientOnly>
            <VantaBackground />
          </ClientOnly>
        )}
        <TanStackQueryProvider>
          <TooltipProvider>{children}</TooltipProvider>
        </TanStackQueryProvider>
        <Scripts />
      </body>
    </html>
  )
}

function RootErrorComponent({ error }: { error: unknown }) {
  const message = error instanceof Error ? error.message : 'Unknown route error'

  return (
    <main className="p-6 space-y-4">
      <h1 className="text-xl font-semibold">Something went wrong</h1>
      <p className="text-sm text-muted-foreground">{message}</p>
    </main>
  )
}
