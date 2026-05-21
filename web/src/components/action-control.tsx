import { NULL_CONTEXT, type Context } from '#/lib/game/context'
import { Button } from './ui/button'
import { useStore } from '@tanstack/react-store'
import { sendContextMessage } from '#/lib/stores/socket'
import { gameStore } from '#/lib/stores/game'
import { clientsStore } from '#/lib/stores/clients'
import { TargetButton } from './target-button'
import type { Action, ActionTransaction } from '#/lib/game/action'
import { setActionID } from '#/lib/stores/battle-context'
import { useValidateContext } from '#/hooks/use-validate-context'
import { useGetTargets } from '#/hooks/use-get-targets'
import { ChevronRight } from 'lucide-react'
import { memo, useMemo } from 'react'

const ActionControl = memo(function ActionControl({
  action,
  staged,
  enabled,
  context,
  onContextChange,
}: {
  action?: Action
  staged?: ActionTransaction
  enabled: boolean
  context: Context
  onContextChange: (context: Context) => void
}) {
  const { valid, loading } = useValidateContext(context)

  const client = useStore(clientsStore, (c) => c.me!)
  const { context: t_context } = useGetTargets(context)

  const actors_all = useStore(gameStore, (g) => g.actors)
  const players_all = useStore(gameStore, (g) => g.players)
  const queued_actions = useStore(gameStore, (g) => g.queued_actions)

  const actors = useMemo(
    () => actors_all.filter((a) => t_context?.target_actor_IDs?.includes(a.ID)),
    [actors_all, t_context?.target_actor_IDs]
  )

  const enemy = useMemo(
    () => players_all.find((p) => p.ID !== client.ID),
    [players_all, client.ID]
  )
  const player = useMemo(
    () => players_all.find((p) => p.ID === client.ID),
    [players_all, client.ID]
  )
  const has_queued_action = useMemo(
    () => queued_actions[context.source_actor_ID ?? ''],
    [queued_actions, context.source_actor_ID]
  )

  if (!!staged) {
    return (
      <div className="flex flex-col py-4 items-center">
        {has_queued_action ? (
          <span className="text-muted-foreground">
            {staged.mutation.config.name}
          </span>
        ) : (
          <Button
            disabled={!enabled}
            onClick={() => {
              sendContextMessage({
                type: 'remove-action',
                client_ID: client.ID,
                context: {
                  ...NULL_CONTEXT,
                  action_ID: staged.ID,
                },
              })
            }}
          >
            Cancel {staged.mutation.config.name}
          </Button>
        )}
      </div>
    )
  }

  return (
    <div className="flex flex-col items-center gap-4 min-w-xs">
      {action &&
        (action.config.switch ? (
          <div className="gap-3 grid grid-cols-2">
            {actors.map((a) => {
              return (
                <TargetButton
                  key={a.ID}
                  actor={a}
                  enabled={enabled}
                  loading={false}
                  contextValid={!!valid}
                  targetType={action.config.target_type}
                  context={context}
                  onContextChange={onContextChange}
                />
              )
            })}
          </div>
        ) : (
          <div className="flex flex-col gap-2">
            {enemy && (
              <div className="gap-3 grid grid-cols-2">
                {enemy.positions.map((pos) => {
                  const a = actors.find((a) => a.ID === pos.actor_ID)
                  if (!a) return <div key={pos.ID} />
                  return (
                    <TargetButton
                      key={a.ID}
                      actor={a}
                      enabled={enabled}
                      loading={false}
                      contextValid={!!valid}
                      targetType={action.config.target_type}
                      context={context}
                      onContextChange={onContextChange}
                    />
                  )
                })}
              </div>
            )}
            {player && (
              <div className="gap-3 grid grid-cols-2">
                {player.positions.map((pos) => {
                  const a = actors.find((a) => a.ID === pos.actor_ID)
                  if (!a) return <div key={pos.ID} />
                  return (
                    <TargetButton
                      key={a.ID}
                      actor={a}
                      enabled={enabled}
                      loading={false}
                      contextValid={!!valid}
                      targetType={action.config.target_type}
                      context={context}
                      onContextChange={onContextChange}
                    />
                  )
                })}
              </div>
            )}
            {actors.length == 0 && valid === false && loading === false && (
              <span className="text-muted-foreground text-sm mb-4">
                no targets available
              </span>
            )}
            {actors.length == 0 && valid === true && loading === false && (
              <span className="text-muted-foreground/50 text-sm mb-4">
                this action does not require selection
              </span>
            )}
          </div>
        ))}

      <div className="flex w-full justify-end">
        <Button
          disabled={!(enabled && action && valid)}
          onClick={() => {
            sendContextMessage({
              type: 'push-action',
              client_ID: client.ID,
              context,
            })

            setActionID(
              context.source_actor_ID!,
              context.action_ID!,
              gameStore.state
            )
          }}
        >
          Next
          <ChevronRight />
        </Button>
      </div>
    </div>
  )
})

export { ActionControl }
