import { getRequest } from '@tanstack/react-start/server'

/**
 * Base URL for the Go API from SSR server functions.
 * Set API_URL (runtime) or VITE_BACKEND_URL (build) in production, or rely on the request origin when /api is on the same host.
 */
export function getApiBaseUrl(): string {
  const fromEnv =
    (typeof process !== 'undefined' && process.env.API_URL) ||
    import.meta.env.VITE_BACKEND_URL
  if (fromEnv && String(fromEnv).trim()) {
    return String(fromEnv).trim().replace(/\/$/, '')
  }

  const req = getRequest()
  if (req?.url) {
    try {
      return new URL(req.url).origin
    } catch {
      /* ignore */
    }
  }
  return ''
}
