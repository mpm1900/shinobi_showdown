import { MODIFIER_ICONS, SHINOBI_ICONS } from '#/data/icons'
import type { Actor } from '#/lib/game/actor'
import { natureIndexes } from '#/lib/game/nature'
import { clientsStore } from '#/lib/stores/clients'
import { cn, keys } from '#/lib/utils'
import { useStore } from '@tanstack/react-store'
import { NatureBadge } from './nature-badge'

function LobbyActorDetails({
  actor,
  className,
  enabled,
  ...props
}: React.ComponentProps<'button'> & { actor: Actor; enabled: boolean }) {
  const client = useStore(clientsStore, (s) => s.me)
  const ItemIcon = actor?.item?.icon
    ? MODIFIER_ICONS[actor.item.icon]
    : undefined
  const AbilityIcon = actor?.ability?.icon
    ? MODIFIER_ICONS[actor.ability.icon]
    : undefined
  return (
    <button
      {...props}
      className={cn(
        'relative overflow-hidden border border-stone-300/20 ring ring-black bg-mauve-900 rounded text-left',
        client?.ID === actor.player_ID && 'bg-slate-900 cursor-pointer',
        enabled && 'bg-slate-600',
        'flex flex-col gap-2',
        className
      )}
    >
      <div className='top-0 bottom-0 left-0 w-1/2 absolute bg-gradient-to-r from-black to-transparent' />
      <img
        src={actor.sprite_url}
        draggable={false}
        className={cn('absolute left-50 bottom-0 opacity-40')}
        style={{
          imageRendering: 'pixelated',
        }}
        width={128}
        height={128}
      />
      <div className="absolute z-0 opacity-20 -right-5 -bottom-7">
        {actor.affiliations
          ?.filter((_, i) => i == 0)
          .map((a) => {
            const C = SHINOBI_ICONS[a]
            return C ? <C key={a} className="w-36" /> : null
          })}
      </div>
      <div
        className={cn(
          'p-2 flex justify-between items-center z-10 border-b border-black',
          client?.ID === actor.player_ID
            ? enabled
              ? 'bg-slate-800'
              : 'bg-slate-950 cursor-pointer'
            : enabled
              ? 'bg-mauve-800'
              : 'bg-mauve-950'
        )}
      >
        <div
          className={cn(
            'text-3xl h-7 nanum-brush-script-regular text-shadow-[2px_2px_0px_#000000]'
          )}
        >
          {actor.name}
        </div>
        <div className="flex items-start">
          {keys(actor.natures)
            .sort((a, b) => natureIndexes[a] - natureIndexes[b])
            .map((nature) => (
              <NatureBadge
                key={nature}
                nature={nature}
                className="block text-xs"
              />
            ))}
        </div>
      </div>
      <div className="grid grid-cols-2 p-2">
        <table className="z-10 [&_td]:px-2 [&_td]:whitespace-nowrap text-shadow-[1px_1px_0px_#000000]">
          <tbody>
            {actor.actions
              .filter((a) => !a.config.switch)
              .map((a) => (
                <tr key={a.ID}>
                  <td className="w-6">
                    {a.config.nature && (
                      <NatureBadge nature={a.config.nature} />
                    )}
                  </td>
                  <td className="font-semibold">{a.config.name}</td>
                </tr>
              ))}
          </tbody>
        </table>
        <div className="p-3 bg-black/60 rounded-xs overflow-hidden ring ring-black mb-2 mr-2 text-shadow-[1px_1px_0px_#000000]">
          <div className="capitalize">{actor.focus}</div>
          <div className="flex gap-1 items-center">
            {ItemIcon && <ItemIcon />}
            <div className="truncate">{actor.item?.name ?? '-'}</div>
          </div>
          <div className="flex gap-1 items-center">
            {AbilityIcon && <AbilityIcon />}
            <div className="truncate">{actor.ability?.name ?? '-'}</div>
          </div>
        </div>
      </div>
    </button>
  )
}

export { LobbyActorDetails }
