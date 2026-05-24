import { getRequest } from '@tanstack/react-start/server'

/**
 * Public origin for the Go API during SSR (same host as the site when /api is reverse-proxied).
 *
 * Prefer `API_URL` on the web container. Behind Traefik/nginx, forwarded headers are used when
 * `getRequest().url` is an internal URL. Server-side fetch requires an absolute URL.
 */
export function getApiBaseUrl(): string {
  const fromEnv =
    (typeof process !== 'undefined' && process.env.API_URL) ||
    import.meta.env.VITE_BACKEND_URL
  if (fromEnv && String(fromEnv).trim()) {
    return String(fromEnv).trim().replace(/\/$/, '')
  }

  const req = getRequest()
  if (!req) {
    throw new Error(
      'API base URL: no request context. Set API_URL (e.g. https://your.domain) on the web container.'
    )
  }

  const xfHost = req.headers.get('x-forwarded-host')
  const xfProto = req.headers.get('x-forwarded-proto')
  if (xfHost) {
    const proto = (xfProto?.split(',')[0].trim() || 'https').replace(/:$/, '')
    try {
      return new URL(`${proto}://${xfHost}`).origin
    } catch {
      /* fall through */
    }
  }

  const host = req.headers.get('host')
  if (host) {
    try {
      const hintedProto = xfProto?.split(',')[0].trim() || 'https'
      return new URL(`${hintedProto}://${host}`).origin
    } catch {
      /* fall through */
    }
  }

  if (req.url) {
    try {
      return new URL(req.url).origin
    } catch {
      /* fall through */
    }
  }

  throw new Error(
    'API base URL: set API_URL (public https origin, no trailing slash) on the web container.'
  )
}
