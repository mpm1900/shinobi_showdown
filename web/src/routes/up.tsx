import { createFileRoute } from '@tanstack/react-router'

/**
 * Kamal/kamal-proxy default health check path (GET /up, expects 2xx).
 */
export const Route = createFileRoute('/up')({
  component: () => 'ok',
})
