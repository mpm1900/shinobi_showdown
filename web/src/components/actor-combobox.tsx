import type { ActorDef } from '#/lib/game/actor'
import { type NatureSet } from '#/lib/game/nature'
import { actorsQuery } from '#/lib/queries/actors'
import { cn } from '#/lib/utils'
import { useSuspenseQuery } from '@tanstack/react-query'
import { Loader2, Plus } from 'lucide-react'
import { useEffect, useMemo, useState, useTransition } from 'react'
import { NatureBadge } from './nature-badge'
import { buttonVariants } from './ui/button'
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
  ComboboxTrigger,
  ComboboxValue,
} from './ui/combobox'
import { Separator } from './ui/separator'

function ActorCombobox({
  className,
  onClick,
  selected = [],
  active,
  value,
  onValueChange,
}: {
  className?: string
  onClick?: () => void
  selected?: string[]
  active?: string
  value?: string | null
  onValueChange?: (def: ActorDef) => void
}) {
  const actors = useSuspenseQuery(actorsQuery)
  const [open, setOpen] = useState(false)
  const [showOpenDetails, setShowOpenDetails] = useState(false)
  const [optimisticActorId, setOptimisticActorId] = useState<string | null>(
    null
  )
  const [isPending, startTransition] = useTransition()
  const sortedActors = useMemo(
    () => [...actors.data].sort((a, b) => a.name.localeCompare(b.name)),
    [actors.data]
  )
  const actorById = useMemo(
    () => new Map(actors.data.map((a) => [a.actor_ID, a])),
    [actors.data]
  )
  const selectedSet = useMemo(() => new Set(selected), [selected])
  const effectiveValueId = optimisticActorId ?? value ?? null
  const actor =
    (effectiveValueId ? actorById.get(effectiveValueId) : null) ?? null
  const is_active = !!actor && active === actor.actor_ID
  const selected_actors = selected.map((id) => actorById.get(id))
  const has_restricted = selected_actors.some((a) => a?.restricted)
  const showRestrictedDisabled = !actor?.restricted && has_restricted
  const disabledActorIds = useMemo(
    () =>
      new Set(
        sortedActors
          .filter(
            (a) =>
              selectedSet.has(a.actor_ID) ||
              (showRestrictedDisabled && a.restricted)
          )
          .map((a) => a.actor_ID)
      ),
    [sortedActors, selectedSet, showRestrictedDisabled]
  )

  useEffect(() => {
    setOptimisticActorId(null)
  }, [value])

  useEffect(() => {
    if (!open) {
      setShowOpenDetails(false)
      return
    }

    const raf = requestAnimationFrame(() => setShowOpenDetails(true))
    return () => cancelAnimationFrame(raf)
  }, [open])

  const handleValueChange = (actor: ActorDef | null) => {
    if (!actor?.actor_ID) return
    setOpen(false)
    setOptimisticActorId(actor.actor_ID)
    startTransition(() => {
      onValueChange?.(actor)
    })
  }
  return (
    <Combobox<ActorDef>
      items={sortedActors}
      itemToStringValue={(actor) => actor.actor_ID}
      itemToStringLabel={(actor) => actor.name}
      value={actor}
      open={open}
      onOpenChange={setOpen}
      onValueChange={handleValueChange}
    >
      <div
        className={cn(
          buttonVariants({ variant: 'outline' }),
          'relative flex justify-start h-auto p-0!',
          {
            'bg-accent! border-stone-300/50!': is_active,
            'border border-amber-400/40!': actor?.restricted,
            'border border-amber-400!': actor?.restricted && is_active,
            'opacity-80': isPending,
          },
          className
        )}
      >
        {actor ? (
          <img
            src={actor.sprite_url}
            className={cn(
              'size-15 p-0.5 mr-0 bg-stone-300/40 border border-stone-950 rounded cursor-pointer select-none',
              is_active && 'bg-stone-300'
            )}
            onPointerDown={(e) => {
              e.preventDefault()
              e.stopPropagation()
            }}
            onClick={(e) => {
              e.preventDefault()
              e.stopPropagation()
              onClick?.()
            }}
          />
        ) : (
          <Plus className="text-muted-foreground/60 size-8 m-2 mr-0" />
        )}
        <ComboboxTrigger
          className={cn(
            'relative mr-2 flex min-w-0 flex-1 items-center justify-between gap-4',
            actor?.restricted && '*:text-amber-400!'
          )}
        >
          <div className="flex min-w-0 flex-1 flex-col items-start justify-start">
            <div
              className={cn(
                'text-md w-full min-w-0 font-semibold',
                !is_active && 'text-muted-foreground!'
              )}
            >
              <ComboboxValue>
                <div className="flex min-w-0 items-center">
                  {!actor && (
                    <div className="block min-w-0 flex-1 truncate py-3">
                      Select a shinobi...
                    </div>
                  )}
                  {actor && (
                    <span
                      className={cn('min-w-0 flex-1 truncate text-left', {
                        'text-amber-400/60!': actor?.restricted,
                        'text-amber-400!': actor?.restricted && is_active,
                      })}
                    >
                      {actor?.name}
                    </span>
                  )}
                  {actor && (
                    <span className="ml-2 shrink-0 text-xs font-black text-stone-300/30">
                      {isPending ? (
                        <span className="inline-flex items-center gap-1 text-stone-300/60">
                          <Loader2 className="size-3 animate-spin" />
                          Updating...
                        </span>
                      ) : (
                        'Lv.100'
                      )}
                    </span>
                  )}
                </div>
              </ComboboxValue>
            </div>
            {actor && <Separator className="mb-1" />}
            {actor && (
              <div
                className={cn('flex items-start', !is_active && 'opacity-50')}
              >
                {(Object.keys(actor.natures) as Array<NatureSet>).map(
                  (nature) => (
                    <NatureBadge
                      key={nature}
                      nature={nature}
                      className="text-xs"
                    />
                  )
                )}
              </div>
            )}
          </div>
        </ComboboxTrigger>

        {open && (
          <ComboboxContent className="min-w-(--anchor-width) w-(--anchor-width) max-w-(--anchor-width)">
            <ComboboxInput showTrigger={false} placeholder="Search" />
            <ComboboxEmpty>No Shinobi found.</ComboboxEmpty>
            <ComboboxList>
              {(a) => (
                <ComboboxItem
                  key={a.actor_ID}
                  value={a}
                  disabled={disabledActorIds.has(a.actor_ID)}
                  className={cn('justify-between', {
                    'text-amber-400': a.restricted,
                  })}
                  showCheck={false}
                >
                  <div className="truncate">{a.name}</div>
                  {showOpenDetails ? (
                    <div className="flex">
                      {(Object.keys(a.natures) as Array<NatureSet>).map(
                        (nature) => (
                          <NatureBadge
                            key={nature}
                            nature={nature}
                            className="text-xs"
                          />
                        )
                      )}
                    </div>
                  ) : (
                    <div className="h-5" />
                  )}
                </ComboboxItem>
              )}
            </ComboboxList>
          </ComboboxContent>
        )}
      </div>
    </Combobox>
  )
}

export { ActorCombobox }
