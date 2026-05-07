import { getTargets, type Context } from '#/lib/game/context'
import { clientsStore } from '#/lib/stores/clients'
import { gameStore } from '#/lib/stores/game'
import { cn } from '#/lib/utils'
import { useStore } from '@tanstack/react-store'

function RunningContext({ context }: { context: Context }) {
  const game = useStore(gameStore, (g) => g)
  const client_ID = useStore(clientsStore, (s) => s.me?.ID)
  const source = game.actors.find((a) => a.ID === context.source_actor_ID)
  const source_action = source?.actions.find((a) => a.ID === context.action_ID)
  const targets = getTargets(source_action?.target_type, game, context)
  const has_targets =
    targets.length > 0 && targets[0].ID !== context.source_actor_ID

  if (!source || !source_action) return null

  const is_friendly_source = client_ID === source.player_ID

  return (
    <div className="pointer-events-none absolute inset-0 z-30 grid place-items-center px-4">
      <div
        key={context.source_actor_ID}
        className="relative w-full max-w-2xl animate-in fade-in-0 zoom-in-95 duration-300"
      >
        <div
          className={cn(
            'absolute -inset-4 -z-10 rounded-2xl opacity-70 blur-2xl',
            {
              'bg-blue-600/25': is_friendly_source,
              'bg-rose-700/25': !is_friendly_source,
            }
          )}
        />

        <div
          className={cn(
            'relative w-full overflow-hidden rounded-lg border bg-stone-900/92 ring ring-black shadow-[0_10px_30px_rgba(0,0,0,0.7)]',
            {
              'border-blue-200/30': is_friendly_source,
              'border-rose-200/30': !is_friendly_source,
            }
          )}
        >
          <div className="px-5 pt-3 pb-3">
            <div className="flex flex-wrap items-center justify-center gap-x-3 gap-y-2 leading-none">
              <span
                className={cn(
                  'text-3xl sm:text-4xl tracking-tight text-shadow-[1px_1px_0px_#000000] nanum-brush-script-regular',
                  {
                    'text-blue-300': is_friendly_source,
                    'text-rose-300': !is_friendly_source,
                  }
                )}
              >
                {source.name}
              </span>
              <span className="pb-1 text-xs font-bold uppercase tracking-[0.2em] text-stone-300/90">
                uses
              </span>
              <span className="text-5xl tracking-wide text-amber-200 text-shadow-[1px_1px_0px_#000000] nanum-brush-script-regular">
                {source_action.config.name}
              </span>
            </div>

            {has_targets && (
              <div className="mt-3 h-px w-full bg-gradient-to-r from-transparent via-stone-100/35 to-transparent" />
            )}

            {has_targets && (
              <div className="mt-2.5 flex flex-wrap items-center justify-center gap-2">
                <span className="text-[10px] font-bold uppercase text-stone-300/60">
                  On
                </span>
                {targets.map((target) => (
                  <span
                    key={target.ID}
                    className="text-sm font-semibold text-stone-100 text-shadow-[1px_1px_0px_#000000]"
                  >
                    {target.name}
                  </span>
                ))}
              </div>
            )}
          </div>

          <div className="h-1 w-full bg-[linear-gradient(90deg,theme(colors.transparent)_0%,theme(colors.amber.300)_50%,theme(colors.transparent)_100%)] opacity-60" />
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_50%_20%,rgba(255,255,255,0.1),transparent_60%)]" />
          <div className="absolute inset-0 bg-[repeating-linear-gradient(-45deg,transparent,transparent_18px,rgba(255,255,255,0.03)_18px,rgba(255,255,255,0.03)_20px)] opacity-40" />
        </div>

        <div className="absolute -bottom-1 left-5 right-50 h-px bg-gradient-to-r from-transparent via-amber-200/60 to-transparent" />
        <div className="absolute -top-1 left-50 right-5 h-px bg-gradient-to-r from-transparent via-amber-200/60 to-transparent" />
      </div>

      <div
        className={cn('absolute -z-20 h-44 w-44 rounded-full blur-3xl', {
          'bg-blue-400/15': is_friendly_source,
          'bg-rose-400/15': !is_friendly_source,
        })}
      />
    </div>
  )
}

export { RunningContext }
