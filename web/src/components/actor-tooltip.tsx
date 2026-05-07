import { useActiveActor } from '#/hooks/use-active-actor'
import { useGetTargets } from '#/hooks/use-get-targets'
import type { Actor } from '#/lib/game/actor'
import { natureIndexes, type NatureSet } from '#/lib/game/nature'
import { ActorStat } from './actor-stat'
import { NatureBadge } from './nature-badge'
import { HoverCard, HoverCardContent, HoverCardTrigger } from './ui/hover-card'
import { NULL_CONTEXT } from '#/lib/game/context'
import { Button } from './ui/button'
import { useStore } from '@tanstack/react-store'
import { gameStore } from '#/lib/stores/game'
import { clientsStore } from '#/lib/stores/clients'
import { sendContextMessage } from '#/lib/stores/socket'
import { setActionID } from '#/lib/stores/battle-context'
import { cn, keys } from '#/lib/utils'
import { Separator } from './ui/separator'
import { MODIFIER_ICONS } from '#/data/icons'

function SwitchButton({ actor }: { actor: Actor }) {
  const active = useActiveActor()
  const switch_action = active?.actions.find((a) => a.meta.switch)
  const context = {
    ...NULL_CONTEXT,
    action_ID: switch_action?.ID ?? null,
    parent_actor_ID: active?.ID ?? null,
    source_actor_ID: active?.ID ?? null,
    source_player_ID: active?.player_ID ?? null,
    target_actor_IDs: [],
  }
  const client = useStore(clientsStore, (c) => c.me!)
  const game = useStore(gameStore, (g) => g)
  const idle = game.status === 'idle'
  const queued_action = game.actions.find(
    (a) =>
      a.context.parent_actor_ID === active?.ID ||
      (a.context.target_actor_IDs?.includes(actor.ID) &&
        a.mutation.ID === switch_action?.ID)
  )

  const { context: t_context } = useGetTargets(context)
  if (
    switch_action?.disabled ||
    !idle ||
    !t_context?.target_actor_IDs?.includes(actor.ID)
  ) {
    return null
  }

  return (
    <Button
      disabled={!!queued_action}
      onClick={() => {
        sendContextMessage({
          type: 'push-action',
          client_ID: client.ID,
          context: {
            ...context,
            target_actor_IDs: [actor.ID],
          },
        })

        setActionID(context.source_actor_ID!, context.action_ID!, game)
      }}
    >
      Switch {actor.name}
    </Button>
  )
}

function ActorTooltip({
  actor,
  disabled = false,
  ...props
}: React.ComponentProps<typeof HoverCardTrigger> & {
  actor: Actor
  disabled?: boolean
}) {
  const ItemIcon = actor?.item?.icon
    ? MODIFIER_ICONS[actor.item.icon]
    : undefined
  const AbilityIcon = actor?.ability?.icon
    ? MODIFIER_ICONS[actor.ability.icon]
    : undefined
  return (
    <HoverCard openDelay={300} closeDelay={100}>
      <HoverCardTrigger {...props} />
      {!disabled && (
        <HoverCardContent
          sideOffset={8}
          collisionPadding={8}
          className="w-auto flex flex-col"
        >
          <div className="flex justify-between">
            <div className="text-3xl nanum-brush-script-regular">
              {actor.name}
            </div>
            <div>
              {keys(actor.natures)
                .sort((a, b) => natureIndexes[a] - natureIndexes[b])
                .map((nature) => (
                  <NatureBadge key={nature} nature={nature} />
                ))}
            </div>
          </div>
          <div className="bg-stone-700 rounded-xs overflow-hidden ring ring-black mb-2">
            <div className="mb-1 h-px w-full bg-gradient-to-r to-stone-100/35 from-transparent" />
            <div className="flex gap-8 [&>div]:flex-1 [&>div]:text-nowrap px-2">
              <div className="capitalize">{actor.focus}</div>
              <div className="flex gap-1 items-center">
                {ItemIcon && <ItemIcon />}
                {actor.item?.name ?? '-'}
              </div>
              <div className="flex gap-1 items-center">
                {AbilityIcon && <AbilityIcon />}
                {actor.ability?.name ?? '-'}
              </div>
            </div>
            <div className="mt-1 h-px w-full bg-gradient-to-r to-transparent from-stone-100/35" />
          </div>

          <div className="flex items-start mb-4">
            <table className="[&_td]:px-2 [&_td]:whitespace-nowrap">
              <tbody>
                <tr>
                  <td>HP</td>
                  <td>
                    <ActorStat actor={actor} showBase={false} stat={'hp'} />
                  </td>
                </tr>
                <tr>
                  <td>Stamina</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'stamina'}
                    />
                  </td>
                </tr>
                <tr>
                  <td>Attack</td>
                  <td>
                    <ActorStat actor={actor} showBase={false} stat={'attack'} />
                  </td>
                </tr>
                <tr>
                  <td>Defense</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'defense'}
                    />
                  </td>
                </tr>
                <tr>
                  <td>C.Attack</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'chakra_attack'}
                    />
                  </td>
                </tr>
                <tr>
                  <td>C.Defense</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'chakra_defense'}
                    />
                  </td>
                </tr>
                <tr>
                  <td>Speed</td>
                  <td>
                    <ActorStat actor={actor} showBase={false} stat={'speed'} />
                  </td>
                </tr>
                <tr>
                  <td>Accuracy</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'accuracy'}
                    />
                  </td>
                </tr>
                <tr>
                  <td>Evasion</td>
                  <td>
                    <ActorStat
                      actor={actor}
                      showBase={false}
                      stat={'evasion'}
                    />
                  </td>
                </tr>
              </tbody>
            </table>
            <Separator orientation="vertical" />
            <table className="[&_td]:px-2 [&_td]:whitespace-nowrap">
              <tbody>
                {actor.actions
                  .filter((a) => !a.meta.switch)
                  .map((a) => (
                    <tr
                      key={a.ID}
                      className={cn({
                        'text-destructive': a.disabled,
                      })}
                    >
                      <td>
                        {a.config.nature && (
                          <NatureBadge nature={a.config.nature} />
                        )}
                      </td>
                      <td>{a.config.name}</td>
                    </tr>
                  ))}
              </tbody>
            </table>
          </div>
          <SwitchButton actor={actor} />
        </HoverCardContent>
      )}
    </HoverCard>
  )
}

export { ActorTooltip }
