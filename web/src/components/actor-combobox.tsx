import type { ActorDef } from '#/lib/game/actor'
import { actorsQuery } from '#/lib/queries/actors'
import { useSuspenseQuery } from '@tanstack/react-query'
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
import { Plus } from 'lucide-react'
import { natureIndexes, type NatureSet } from '#/lib/game/nature'
import { NatureBadge } from './nature-badge'
import { buttonVariants } from './ui/button'
import { cn } from '#/lib/utils'
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
  const sortedActors = actors.data.sort((a, b) => a.name.localeCompare(b.name))
  const actor = sortedActors.find((actor) => actor.actor_ID === value) ?? null
  const is_active = !!actor && active === actor.actor_ID
  const selected_actors = selected.map((id) =>
    actors.data.find((a) => a.actor_ID === id)
  )
  const has_restricted = selected_actors.some((a) => a?.restricted)

  const handleValueChange = (actor: ActorDef | null) => {
    if (!actor?.actor_ID) return
    onValueChange?.(actor)
  }
  return (
    <Combobox<ActorDef>
      items={sortedActors}
      itemToStringValue={(actor) => actor.actor_ID}
      itemToStringLabel={(actor) => actor.name}
      value={actor}
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
            'relative flex flex-1 justify-between w-full mr-2 gap-4',
            actor?.restricted && '*:text-amber-400!'
          )}
        >
          <div className="flex flex-col justify-start flex-1 items-start">
            <div
              className={cn(
                'font-semibold text-md',
                !is_active && 'text-muted-foreground!'
              )}
            >
              <ComboboxValue>
                {!actor && <div className="py-3">Select a shinobi...</div>}
                {actor && (
                  <span
                    className={cn({
                      'text-amber-400/60!': actor?.restricted,
                      'text-amber-400!': actor?.restricted && is_active,
                    })}
                  >
                    {actor?.name}
                  </span>
                )}
                {actor && (
                  <span className="text-stone-300/30 text-xs font-black ml-2">
                    Lv.100
                  </span>
                )}
              </ComboboxValue>
            </div>
            {actor && <Separator className="mb-1" />}
            {actor && (
              <div
                className={cn('flex items-start', !is_active && 'opacity-50')}
              >
                {(Object.keys(actor.natures) as Array<NatureSet>)
                  .sort((a, b) => natureIndexes[a] - natureIndexes[b])
                  .map((nature) => (
                    <NatureBadge
                      key={nature}
                      nature={nature}
                      className="text-xs"
                    />
                  ))}
              </div>
            )}
          </div>
        </ComboboxTrigger>

        <ComboboxContent className="min-w-(--anchor-width) w-(--anchor-width) max-w-(--anchor-width)">
          <ComboboxInput showTrigger={false} placeholder="Search" />
          <ComboboxEmpty>No Shinobi found.</ComboboxEmpty>
          <ComboboxList>
            {(a) => (
              <ComboboxItem
                key={a.actor_ID}
                value={a}
                disabled={
                  selected.includes(a.actor_ID) ||
                  (!actor?.restricted && a.restricted && has_restricted)
                }
                className={cn('justify-between', {
                  'text-amber-400': a.restricted,
                })}
              >
                <div>{a.name}</div>
                <div>
                  {(Object.keys(a.natures) as Array<NatureSet>)
                    .sort((a, b) => natureIndexes[a] - natureIndexes[b])
                    .map((nature) => (
                      <NatureBadge
                        key={nature}
                        nature={nature}
                        className="text-xs"
                      />
                    ))}
                </div>
              </ComboboxItem>
            )}
          </ComboboxList>
        </ComboboxContent>
      </div>
    </Combobox>
  )
}

export { ActorCombobox }
