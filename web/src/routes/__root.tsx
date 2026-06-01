import {
  ClientOnly,
  HeadContent,
  Scripts,
  createRootRouteWithContext,
} from '@tanstack/react-router'
import TanStackQueryProvider from '../integrations/tanstack-query/root-provider'

import styles from '../styles.css?url'

import { TooltipProvider } from '#/components/ui/tooltip'
import { VantaBackground } from '#/components/vanta-background'
import type { User } from '#/lib/queries/auth'
import { Toaster } from '@/components/ui/sonner'
import type { QueryClient } from '@tanstack/react-query'

interface RouterContext {
  queryClient: QueryClient
  auth: {
    user: User | null
  }
}

// const THEME_INIT_SCRIPT = `(function(){try{var stored=window.localStorage.getItem('theme');var mode=(stored==='light'||stored==='dark'||stored==='auto')?stored:'auto';var prefersDark=window.matchMedia('(prefers-color-scheme: dark)').matches;var resolved=mode==='auto'?(prefersDark?'dark':'light'):mode;var root=document.documentElement;root.classList.remove('light','dark');root.classList.add(resolved);if(mode==='auto'){root.removeAttribute('data-theme')}else{root.setAttribute('data-theme',mode)}root.style.colorScheme=resolved;}catch(e){}})();`

import { meQuery } from '#/lib/queries/auth'
import { uiStore } from '#/lib/stores/ui'
import { useStore } from '@tanstack/react-store'
import { Login } from './login'

export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async ({ context, location }) => {
    if (location.pathname === '/up') {
      return { auth: { user: null } } as any
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
        content:
          'width=device-width, initial-scale=0.95, maximum-scale=1.0, user-scalable=no, viewport-fit=cover',
      },
      {
        name: 'apple-mobile-web-app-capable',
        content: 'yes',
      },
      {
        name: 'apple-mobile-web-app-status-bar-style',
        content: 'black-translucent',
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
  errorComponent: Login,
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
        <Toaster position="bottom-center" />
        <Scripts />
      </body>
    </html>
  )
}
