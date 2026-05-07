import type { Actor } from '#/lib/game/actor'
import { cn } from '#/lib/utils'

function ActorStatus({ actor }: { actor: Actor }) {
  return (
    <div
      className={cn(
        'absolute font-bold px-1 mx-1 left-1 right-1 h-4 text-center leading-5! rounded-sm whitespace-nowrap -bottom-1 z-10 text-lg nanum-brush-script-regular',
        'bg-stone-200 border border-stone-900 text-stone-900 shadow-[0px_1px_2px_rgba(0,0,0,1)]',
        {
          'bg-orange-200 border-orange-700 text-orange-900': actor.burned,
          'bg-indigo-300 border-indigo-700 text-indigo-900': actor.sleeping,
          'bg-yellow-200 border-yellow-700 text-yellow-900': actor.paralyzed,
          'bg-lime-200 border-lime-700 text-yellow-900': actor.poisoned,
        }
      )}
    >
      {actor.statused ? (
        <span className='-mt-0.5 block'>
          {actor.sleeping && 'SLEEP'}
          {actor.paralyzed && 'PARA'}
          {actor.burned && 'BURN'}
          {actor.poisoned && 'POISON'}
        </span>
      ) : (
        <span className='-mt-0.5 block'>LV {actor.level}</span>
      )}
    </div>
  )
}

export { ActorStatus }
