import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import { useQuery } from '@tanstack/react-query'
import { TeamBuilderActionsTable } from './team-builder-actions-table'
import { TeamBuilderActorConfig } from './team-builder-actor-config'
import { actorsQuery } from '#/lib/queries/actors'

function TeamBuilderActor({ form }: { form: TeamBuilderForm }) {
  const actors = useQuery(actorsQuery)
  return (
    <form.Subscribe
      selector={(state) => ({
        actor: state.values.actors[state.values.selected_index],
      })}
    >
      {({ actor }) => {
        if (!actor) {
          return (
            <div className="flex-1 text-sm text-muted-foreground">
              Select a shinobi portrait to edit config
            </div>
          )
        }

        const def = actors.data?.find((a) => a.actor_ID === actor.actor_ID)

        if (!def) return null

        return (
          <div className="flex-1 flex flex-col gap-2 overflow-auto">
            <TeamBuilderActorConfig form={form} def={def} />
            <TeamBuilderActionsTable form={form} def={def} />
          </div>
        )
      }}
    </form.Subscribe>
  )
}

export { TeamBuilderActor }
