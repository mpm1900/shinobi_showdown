import type { Actor } from '#/lib/game/actor'
import type { Modifier as ModifierType } from '#/lib/game/modifier'
import { Modifier } from './modifier'



function ActorModifiers({
  actor,
  modifiers,
}: {
  actor: Actor
  modifiers: ModifierType[]
}) {
  return (
    <div className="relative flex flex-row-reverse justify-end items-end flex-wrap px-2 gap-2 z-30">
      {Object.entries(actor.applied_modifiers ?? {}).map(([ID, count]) => (
        <Modifier
          key={ID}
          count={count}
          modifier={modifiers.find((m) => m.group_ID === ID)}
        />
      ))}
    </div>
  )
}

export { ActorModifiers }
