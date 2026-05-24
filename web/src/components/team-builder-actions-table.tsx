import type { TeamBuilderForm } from '#/hooks/use-team-builder-form'
import type { ActorDef } from '#/lib/game/actor'
import { actionsQuery } from '#/lib/queries/actions'
import { keys } from '#/lib/utils'
import { useSuspenseQuery } from '@tanstack/react-query'
import { ActionsTable } from './actions-table'

function TeamBuilderActionsTable({
  def,
  form,
}: {
  def: ActorDef | undefined
  form: TeamBuilderForm
}) {
  const actions = useSuspenseQuery(actionsQuery)
  const data = !def?.action_IDs?.length
    ? actions.data
    : actions.data.filter((a) => def?.action_IDs.includes(a.ID))

  return (
    <form.Subscribe
      selector={(state) => ({
        selected_index: state.values.selected_index,
        actors: state.values.actors,
      })}
    >
      {({ selected_index, actors }) => (
        <ActionsTable
          total={def?.action_count ?? 0}
          data={data}
          rowSelection={Object.fromEntries(
            actors[selected_index]?.config.action_IDs?.map((id) => [
              id,
              true,
            ]) ?? []
          )}
          onRowSelectionChange={(selection) => {
            form.setFieldValue(
              `actors[${selected_index}].config.action_IDs`,
              keys(selection)
            )
          }}
        />
      )}
    </form.Subscribe>
  )
}

export { TeamBuilderActionsTable }
